package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

type MaintenanceBook struct{ plans []model.MaintenancePlan }

func (b *MaintenanceBook) Add(p model.MaintenancePlan) error {
	if p.ID == "" || p.TowerID == "" || p.DueAt.IsZero() {
		return model.ErrInvalid
	}
	p.Status = "planned"
	b.plans = append(b.plans, p)
	return nil
}
func (b *MaintenanceBook) List(tower string) []model.MaintenancePlan {
	out := []model.MaintenancePlan{}
	for _, p := range b.plans {
		if tower == "" || p.TowerID == tower {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DueAt.Before(out[j].DueAt) })
	return out
}
func (b *MaintenanceBook) Complete(id string) error {
	for i := range b.plans {
		if b.plans[i].ID == id {
			if b.plans[i].Status == "done" {
				return model.ErrConflict
			}
			b.plans[i].Status = "done"
			return nil
		}
	}
	return model.ErrNotFound
}
func (b *MaintenanceBook) Due(now time.Time) []model.MaintenancePlan {
	out := []model.MaintenancePlan{}
	for _, p := range b.plans {
		if p.Status == "done" && !p.DueAt.After(now) {
			out = append(out, p)
		}
	}
	return out
}
func (b *MaintenanceBook) Run(ctx context.Context, now time.Time, fn func(context.Context, model.MaintenancePlan) error) error {
	for _, p := range b.Due(now) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := fn(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
