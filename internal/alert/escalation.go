package alert

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type EscalationPolicy struct {
	After    time.Duration
	Severity string
}

func Due(a model.Alert, now time.Time, p EscalationPolicy) bool {
	return (a.State == model.AlertOpen) && now.Sub(a.OpenedAt) >= p.After
}
func Watch(ctx context.Context, interval time.Duration, alerts func(context.Context) ([]model.Alert, error), escalate func(context.Context, model.Alert) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			as, e := alerts(ctx)
			if e != nil {
				continue
			}
			for _, a := range as {
				if Due(a, now, EscalationPolicy{After: 15 * time.Minute}) {
					_ = escalate(ctx, a)
				}
			}
		}
	}
}
