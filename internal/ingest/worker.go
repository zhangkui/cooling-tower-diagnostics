package ingest

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
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
	for _, r := range frames {
		if err := submit(ctx, r); err != nil {
			return err
		}
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}
