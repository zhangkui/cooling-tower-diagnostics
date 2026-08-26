package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/ingest"
	"errors"
	"testing"
	"time"
)

func TestBug10_RetryBackoffHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- ingest.Retry(ctx, ingest.RetryPolicy{Attempts: 2, Base: time.Second}, func(context.Context) error { return errors.New("gateway unavailable") }) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) { t.Fatalf("retry error = %v", err) }
	case <-time.After(150 * time.Millisecond):
		t.Fatal("retry remained blocked during cancellation")
	}
}
