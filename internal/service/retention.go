package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type Retention struct{ Store store.RetentionStore }

func (r Retention) Run(ctx context.Context, p model.RetentionPolicy) (int64, error) {
	return r.Store.Purge(ctx, p, time.Now().UTC())
}
func NextRun(now time.Time, hour int) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
