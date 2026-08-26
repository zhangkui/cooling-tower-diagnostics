package clock

import (
	"context"
	"time"
)

type Ticker struct {
	C    <-chan time.Time
	stop func()
}

func NewTicker(d time.Duration) Ticker { t := time.NewTicker(d); return Ticker{C: t.C, stop: t.Stop} }
func (t Ticker) Stop() {
	if t.stop != nil {
		t.stop()
	}
}
func Every(ctx context.Context, d time.Duration, fn func(time.Time)) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			fn(now)
		}
	}
}
