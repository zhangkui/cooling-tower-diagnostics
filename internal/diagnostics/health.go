package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type HealthChecker struct {
	Towers   store.TowerStore
	Readings store.ReadingStore
}

func (h HealthChecker) Check(ctx context.Context, towerID string, now time.Time) (string, error) {
	t, e := h.Towers.Get(ctx, towerID)
	if e != nil {
		return "unknown", e
	}
	r, e := h.Readings.Last(ctx, towerID)
	if e != nil {
		return "no-data", nil
	}
	if now.Sub(r.RecordedAt) > 30*time.Minute {
		return "stale", nil
	}
	if t.Status != model.TowerOnline {
		return "paused", nil
	}
	return "healthy", nil
}
