package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type CapacityService struct {
	Readings store.ReadingStore
	Towers   store.TowerStore
}

func (c CapacityService) Evaluate(ctx context.Context, tower string, rated float64) (model.Capacity, error) {
	if err := ctx.Err(); err != nil {
		return model.Capacity{}, err
	}
	rows, e := c.Readings.Recent(ctx, tower, "flow", time.Now().UTC().Add(-15*time.Minute))
	if e != nil {
		return model.Capacity{}, e
	}
	current := 0.0
	if len(rows) > 0 {
		current = rows[len(rows)-1].Value
	}
	return rules.CapacityStatus(tower, rated, current), nil
}
func (c CapacityService) Trend(ctx context.Context, tower string, rated float64) ([]model.Capacity, error) {
	rows, e := c.Readings.Recent(ctx, tower, "flow", time.Now().UTC().Add(-time.Hour))
	if e != nil {
		return nil, e
	}
	out := []model.Capacity{}
	for _, r := range rows {
		out = append(out, rules.CapacityStatus(tower, rated, r.Value))
	}
	return out, nil
}
