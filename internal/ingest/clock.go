package ingest

import (
	"context"
	"time"
)

type Clock interface{ Now() time.Time }
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC() }
func Await(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func Since(c Clock, at time.Time) time.Duration {
	if c == nil {
		return time.Since(at)
	}
	return c.Now().Sub(at)
}
