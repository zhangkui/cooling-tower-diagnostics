package regression

import (
	"context"
	"cooling-tower-diagnostics/internal/ingest"
	"cooling-tower-diagnostics/internal/model"
	"testing"
	"time"
)

func TestBug01_ReplayStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	submitted := 0
	err := make(chan error, 1)
	started := make(chan struct{})
	go func() {
		err <- ingest.Replay(ctx, []model.ReadingRequest{{Sensor: "a"}, {Sensor: "b"}}, 40*time.Millisecond, func(context.Context, model.ReadingRequest) error {
			submitted++
			if submitted == 1 {
				close(started)
			}
			return nil
		})
	}()
	<-started
	cancel()
	if e := <-err; e == nil {
		t.Fatal("expected cancellation")
	}
	if submitted != 1 {
		t.Fatalf("submitted %d frames after cancellation", submitted)
	}
}
