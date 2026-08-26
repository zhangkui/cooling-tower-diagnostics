package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/service"
	"cooling-tower-diagnostics/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestBug04_UnrelatedEnabledSensorDoesNotCalibrate(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "sensor.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL.ExecContext(ctx, `INSERT INTO towers(id,name,site,design_kw,status,created_at) VALUES(?,?,?,?,?,?)`, "tower-1", "Tower", "Plant", 100, "online", "2026-08-26T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	sensors := store.SensorStore{DB: db}
	if err := sensors.Put(ctx, model.Sensor{ID: "s-1", TowerID: "tower-1", Kind: "conductivity", Scale: 2, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	svc := &service.Service{Clock: fixedClock{at: parseTime("2026-08-26T01:00:00Z")}, Sensors: sensors, Readings: store.ReadingStore{DB: db}}
	reading, err := svc.IngestReading(ctx, model.ReadingRequest{TowerID: "tower-1", Sensor: "outlet_temp", Value: 30, Unit: "C"})
	if err != nil {
		t.Fatal(err)
	}
	if reading.Value != 30 {
		t.Fatalf("value=%v", reading.Value)
	}
}
func parseTime(v string) time.Time { t, _ := time.Parse(time.RFC3339, v); return t }

type fixedClock struct{ at time.Time }

func (c fixedClock) Now() time.Time { return c.at }
