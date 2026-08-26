package ingest

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"errors"
	"sync"
	"time"
)

type Processor func(context.Context, model.ReadingRequest) error
type Pool struct {
	workers int
	queue   chan model.ReadingRequest
	process Processor
	wg      sync.WaitGroup
}

func NewPool(workers, buffer int, p Processor) *Pool {
	if workers < 1 {
		workers = 1
	}
	return &Pool{workers: workers, queue: make(chan model.ReadingRequest, buffer), process: p}
}
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case r, ok := <-p.queue:
					if !ok {
						return
					}
					_ = p.process(ctx, r)
				}
			}
		}()
	}
}
func (p *Pool) Submit(ctx context.Context, r model.ReadingRequest) error {
	select {
	case p.queue <- r:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (p *Pool) Stop() { close(p.queue); p.wg.Wait() }
func Replay(ctx context.Context, frames []model.ReadingRequest, delay time.Duration, submit func(context.Context, model.ReadingRequest) error) error {
	for i, r := range frames {
		// Honor cancellation before submitting each frame so a cancel issued
		// during the inter-frame delay stops playback immediately instead of
		// committing the next frame.
		if err := ctx.Err(); err != nil {
			return errors.Join(model.ErrCanceled, err)
		}
		if err := submit(ctx, r); err != nil {
			return err
		}
		// Only the last frame skips the wait; Await is cancellable so an
		// operator cancel during the delay returns at once.
		if delay > 0 && i < len(frames)-1 {
			if err := Await(ctx, delay); err != nil {
				return errors.Join(model.ErrCanceled, err)
			}
		}
	}
	return nil
}
