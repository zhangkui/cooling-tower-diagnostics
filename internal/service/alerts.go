package service

import (
	"context"
	"cooling-tower-diagnostics/internal/alert"
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

type AlertManager struct {
	Store interface {
		List(context.Context, string) ([]model.Alert, error)
		Get(context.Context, string) (model.Alert, error)
		Upsert(context.Context, model.Alert) error
	}
	Clock func() time.Time
}

func (m AlertManager) Active(ctx context.Context, tower string) ([]model.Alert, error) {
	items, err := m.Store.List(ctx, "")
	if err != nil {
		return nil, err
	}
	out := make([]model.Alert, 0, len(items))
	for _, item := range items {
		if item.TowerID == tower && item.State != model.AlertClosed {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out, nil
}

func (m AlertManager) EscalateDue(ctx context.Context, now time.Time, after time.Duration) (int, error) {
	items, err := m.Store.List(ctx, "")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if item.State == model.AlertAcknowledged && alert.Due(item, now, alert.EscalationPolicy{After: after}) {
			if err := alert.Escalate(&item, now); err != nil {
				continue
			}
			if err := m.Store.Upsert(ctx, item); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

func (m AlertManager) OpenOrSuppress(ctx context.Context, tower, rule, severity, message string) (model.Alert, bool, error) {
	items, err := m.Active(ctx, tower)
	if err != nil {
		return model.Alert{}, false, err
	}
	if alert.ShouldSuppress(items, rule) {
		return model.Alert{}, true, nil
	}
	now := time.Now().UTC()
	if m.Clock != nil {
		now = m.Clock()
	}
	item := alert.Open(now, tower, rule, severity, message)
	if err := m.Store.Upsert(ctx, item); err != nil {
		return item, false, err
	}
	return item, false, nil
}
