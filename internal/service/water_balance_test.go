package service

import (
	"context"
	"cooling-tower-diagnostics/internal/clock"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/ingest"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func testBalanceService(t *testing.T) (*Service, context.Context, time.Time) {
	t.Helper()
	ctx := context.Background()
	root := t.TempDir()
	db, err := store.Open(filepath.Join(root, "cooling.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if err := db.EnsureOptionalSchema(ctx); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 26, 8, 0, 0, 0, time.UTC)
	towers := store.TowerStore{DB: db}
	readings := store.ReadingStore{DB: db}
	svc := &Service{
		Towers: towers, Readings: readings, Audits: store.AuditStore{DB: db},
		Alerts: store.AlertStore{DB: db}, Queries: store.QueryStore{DB: db},
		Devices: store.DeviceStore{DB: db}, Sensors: store.SensorStore{DB: db},
		Maintenance: store.MaintenanceStore{DB: db}, Thresholds: store.ThresholdStore{DB: db},
		Clock: clock.Fixed{Value: now}, Decoder: ingest.Decoder{},
		Engine: diagnostics.Engine{Towers: towers, Readings: readings, Thresholds: store.ThresholdStore{DB: db}, Clock: clock.Fixed{Value: now}},
	}
	return svc, ctx, now
}

func TestAssessWaterBalanceUsesDurableReadings(t *testing.T) {
	svc, ctx, now := testBalanceService(t)
	tower, err := svc.CreateTower(ctx, model.CreateTowerRequest{Name: "Balance Tower", Site: "West", DesignKW: 100})
	if err != nil {
		t.Fatal(err)
	}
	for _, sample := range []struct {
		sensor string
		first  float64
		last   float64
	}{
		{"makeup_flow", 60, 60},
		{"blowdown_flow", 12, 12},
		{"evaporation_flow", 18, 18},
	} {
		for _, at := range []time.Time{now.Add(-10 * time.Minute), now} {
			_, err := svc.IngestReading(ctx, model.ReadingRequest{TowerID: tower.ID, Sensor: sample.sensor, Value: sample.first, Unit: "m3/h", RecordedAt: &at})
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	assessment, err := svc.AssessWaterBalance(ctx, model.WaterBalanceRequest{TowerID: tower.ID, From: now.Add(-10 * time.Minute), To: now, Tolerance: 0.01, MaxGapMinutes: 15})
	if err != nil {
		t.Fatal(err)
	}
	if assessment.Status != "excess-makeup" || assessment.Balance.Makeup != 10 || assessment.Balance.Blowdown != 2 || assessment.Balance.Evaporation != 3 || assessment.Balance.Unaccounted != 5 {
		t.Fatalf("unexpected assessment %#v", assessment)
	}
	audits, err := svc.ListAudit(ctx, "water-balance")
	if err != nil || len(audits) != 1 || audits[0].Action != "assessed" {
		t.Fatalf("expected durable audit, audits=%#v err=%v", audits, err)
	}
}

func TestAssessWaterBalanceReportsMeterGap(t *testing.T) {
	svc, ctx, now := testBalanceService(t)
	tower, err := svc.CreateTower(ctx, model.CreateTowerRequest{Name: "Gap Tower", Site: "East", DesignKW: 100})
	if err != nil {
		t.Fatal(err)
	}
	old := now.Add(-time.Hour)
	_, err = svc.IngestReading(ctx, model.ReadingRequest{TowerID: tower.ID, Sensor: "makeup_flow", Value: 60, Unit: "m3/h", RecordedAt: &old})
	if err != nil {
		t.Fatal(err)
	}
	assessment, err := svc.AssessWaterBalance(ctx, model.WaterBalanceRequest{TowerID: tower.ID, From: old, To: now, MaxGapMinutes: 15})
	if err != nil {
		t.Fatal(err)
	}
	if assessment.Coverage[0].GapCount != 0 || assessment.Coverage[0].Samples != 1 {
		t.Fatalf("single sample should be reported as insufficient, got %#v", assessment.Coverage[0])
	}
	if len(assessment.Findings) == 0 {
		t.Fatal("expected a finding for missing meter samples")
	}
}
