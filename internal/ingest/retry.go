package ingest

import (
	"context"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	Attempts int
	Base     time.Duration
	Jitter   time.Duration
}

func Retry(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	var err error
	for i := 0; i < p.Attempts; i++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		delay := p.Base * time.Duration(i+1)
		if p.Jitter > 0 {
			delay += time.Duration(rand.Int63n(int64(p.Jitter)))
		}
		timer := time.NewTimer(delay)
		<-timer.C
	}
	return err
}
func Backoff(p RetryPolicy, n int) time.Duration {
	if n < 1 {
		n = 1
	}
	return p.Base * time.Duration(n)
}
