package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestBug06_ReportRangeIncludesStartBoundary(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "report.db"))
	if err != nil { t.Fatal(err) }
	defer db.Close()
	if err := db.Migrate(ctx); err != nil { t.Fatal(err) }
	start := time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC)
	if _, err := db.SQL.ExecContext(ctx, `INSERT INTO towers(id,name,site,design_kw,status,created_at) VALUES(?,?,?,?,?,?)`, "tower-1", "Tower", "Plant", 100, "online", start.Format(time.RFC3339Nano)); err != nil { t.Fatal(err) }
	if err := (store.ReadingStore{DB: db}).Add(ctx, model.Reading{ID:"r-start", TowerID:"tower-1", Sensor:"conductivity", Value:1200, RawValue:1200, Unit:"uS/cm", Quality:"good", RecordedAt:start}); err != nil { t.Fatal(err) }
	rows, err := (store.QueryStore{DB: db}).ReadingsByRange(ctx, model.ReportFilter{TowerID:"tower-1", From:start, To:start.Add(time.Hour)})
	if err != nil { t.Fatal(err) }
	if len(rows) != 1 || rows[0].ID != "r-start" { t.Fatalf("range rows=%+v", rows) }
}
