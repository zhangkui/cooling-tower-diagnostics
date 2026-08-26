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

func TestBug09_HeartbeatFreshnessKeepsFreshDeviceOnlineAndMarksOverdueDeviceStale(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "health.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if err := db.EnsureOptionalSchema(ctx); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 26, 11, 0, 0, 0, time.UTC)
	if _, err := db.SQL.ExecContext(ctx, `INSERT INTO towers(id,name,site,design_kw,status,created_at) VALUES(?,?,?,?,?,?)`, "tower-1", "Tower", "Plant", 100, "online", now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	devices := store.DeviceStore{DB: db}
	for _, d := range []model.Device{
		{ID: "fresh", TowerID: "tower-1", Protocol: "modbus", Address: "10.0.0.1", State: "connected", LastHeartbeat: now.Add(-4 * time.Minute)},
		{ID: "overdue", TowerID: "tower-1", Protocol: "modbus", Address: "10.0.0.2", State: "connected", LastHeartbeat: now.Add(-6 * time.Minute)},
	} {
		if err := devices.Put(ctx, d); err != nil {
			t.Fatal(err)
		}
	}
	svc := &service.Service{Clock: fixedClock{now}, Devices: devices, Towers: store.TowerStore{DB: db}, Alerts: store.AlertStore{DB: db}, Audits: store.AuditStore{DB: db}}
	got, err := svc.ReconcileDeviceHealth(ctx, "tower-1")
	if err != nil {
		t.Fatal(err)
	}
	states := map[string]string{}
	for _, d := range got {
		states[d.ID] = d.State
	}
	if states["fresh"] != "connected" || states["overdue"] != "stale" {
		t.Fatalf("states=%v", states)
	}
}

type fixedClock struct{ at time.Time }

func (c fixedClock) Now() time.Time { return c.at }
