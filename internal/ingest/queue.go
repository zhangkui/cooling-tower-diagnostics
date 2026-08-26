package ingest

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"errors"
	"sync"
	"time"
)

var ErrQueueClosed = errors.New("ingest queue closed")

type Queue struct {
	mu     sync.Mutex
	items  []model.ReadingRequest
	closed bool
	notify chan struct{}
}

func NewQueue() *Queue { return &Queue{notify: make(chan struct{}, 1)} }
func (q *Queue) Push(r model.ReadingRequest) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrQueueClosed
	}
	q.items = append(q.items, r)
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return nil
}
func (q *Queue) Pop(ctx context.Context) (model.ReadingRequest, error) {
	for {
		q.mu.Lock()
		if len(q.items) > 0 {
			r := q.items[0]
			q.items = q.items[1:]
			q.mu.Unlock()
			return r, nil
		}
		closed := q.closed
		q.mu.Unlock()
		if closed {
			return model.ReadingRequest{}, ErrQueueClosed
		}
		select {
		case <-ctx.Done():
			return model.ReadingRequest{}, ctx.Err()
		case <-q.notify:
		}
	}
}
func (q *Queue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	select {
	case q.notify <- struct{}{}:
	default:
	}
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
func Drain(ctx context.Context, q *Queue, fn func(model.ReadingRequest) error) error {
	for {
		r, e := q.Pop(ctx)
		if e == ErrQueueClosed {
			return nil
		}
		if e != nil {
			return e
		}
		if e := fn(r); e != nil {
			return e
		}
	}
}
func WaitForQueue(ctx context.Context, q *Queue, max time.Duration) bool {
	deadline := time.NewTimer(max)
	defer deadline.Stop()
	for q.Len() > 0 {
		select {
		case <-ctx.Done():
			return false
		case <-deadline.C:
			return false
		default:
			time.Sleep(time.Millisecond)
		}
	}
	return true
}
