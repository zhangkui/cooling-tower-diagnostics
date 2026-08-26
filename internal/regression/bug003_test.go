package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/alert"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestBug03_OpenAlertCanBeReadBeforeAcknowledgement(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "alerts.db"))
	if err != nil { t.Fatal(err) }
	defer db.Close()
	if err := db.Migrate(ctx); err != nil { t.Fatal(err) }

	opened := alert.Open(time.Now().UTC(), "tower-south", "fan-vibration", model.SeverityWarning, "vibration trend")
	alerts := store.AlertStore{DB: db}
	if err := alerts.Upsert(ctx, opened); err != nil { t.Fatal(err) }
	got, err := alerts.Get(ctx, opened.ID)
	if err != nil { t.Fatalf("get newly opened alert: %v", err) }
	if got.AcknowledgedAt != nil || got.ClosedAt != nil { t.Fatal("fresh alert must not have lifecycle timestamps") }
}
