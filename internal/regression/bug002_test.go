package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/alert"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/service"
	"cooling-tower-diagnostics/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestBug02_DueOpenAlertIsPersistedAsEscalated(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "alerts.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		t.Fatal(err)
	}

	alerts := store.AlertStore{DB: db}
	now := time.Date(2026, 8, 26, 8, 0, 0, 0, time.UTC)
	opened := alert.Open(now.Add(-20*time.Minute), "tower-east", "high-conductivity", model.SeverityCritical, "conductivity is high")
	if _, err := db.SQL.ExecContext(ctx, `INSERT INTO alerts(id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, opened.ID, opened.TowerID, opened.Rule, opened.Severity, opened.State, opened.Message, opened.OpenedAt.Format(time.RFC3339Nano), opened.UpdatedAt.Format(time.RFC3339Nano), "", ""); err != nil {
		t.Fatal(err)
	}

	manager := service.AlertManager{Store: alerts}
	count, err := manager.EscalateDue(ctx, now, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("escalated %d alerts, want 1", count)
	}
	var state string
	if err := db.SQL.QueryRowContext(ctx, "SELECT state FROM alerts WHERE id=?", opened.ID).Scan(&state); err != nil {
		t.Fatal(err)
	}
	if state != model.AlertEscalated {
		t.Fatalf("persisted state = %q, want %q", state, model.AlertEscalated)
	}
}
