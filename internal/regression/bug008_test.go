package regression

import (
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/service"
	"testing"
	"time"
)

func TestBug08_DueReturnsPlannedOverdueWork(t *testing.T) {
	now := time.Date(2026, 8, 26, 10, 0, 0, 0, time.UTC)
	book := service.MaintenanceBook{}
	if err := book.Add(model.MaintenancePlan{ID:"m-1", TowerID:"tower-1", Kind:"lubrication", DueAt:now.Add(-time.Hour)}); err != nil { t.Fatal(err) }
	if got := book.Due(now); len(got) != 1 || got[0].ID != "m-1" { t.Fatalf("due plans=%+v", got) }
}
