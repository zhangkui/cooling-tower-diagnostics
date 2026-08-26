package ingest

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"io"
	"sync"
	"time"
)

type Connector interface {
	Open(context.Context) error
	Close() error
	Read(context.Context) (model.TelemetryFrame, error)
	Name() string
}
type ReplayConnector struct {
	mu     sync.Mutex
	frames []model.TelemetryFrame
	index  int
	opened bool
}

func NewReplayConnector(frames []model.TelemetryFrame) *ReplayConnector {
	return &ReplayConnector{frames: append([]model.TelemetryFrame(nil), frames...)}
}
func (c *ReplayConnector) Open(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opened = true
	c.index = 0
	return nil
}
func (c *ReplayConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opened = false
	return nil
}
func (c *ReplayConnector) Read(ctx context.Context) (model.TelemetryFrame, error) {
	select {
	case <-ctx.Done():
		return model.TelemetryFrame{}, ctx.Err()
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.opened {
		return model.TelemetryFrame{}, fmt.Errorf("connector closed")
	}
	if c.index >= len(c.frames) {
		return model.TelemetryFrame{}, io.EOF
	}
	f := c.frames[c.index]
	c.index++
	return f, nil
}
func (c *ReplayConnector) Name() string { return "replay" }

type ConnectorRunner struct {
	Connector Connector
	Handle    func(context.Context, model.TelemetryFrame) error
	Clock     Clock
}

func (r ConnectorRunner) Run(ctx context.Context) error {
	if r.Connector == nil || r.Handle == nil {
		return model.ErrInvalid
	}
	if err := r.Connector.Open(ctx); err != nil {
		return err
	}
	defer r.Connector.Close()
	for {
		f, e := r.Connector.Read(ctx)
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}
		if err := r.Handle(ctx, f); err != nil {
			return err
		}
	}
}
func Poll(ctx context.Context, connector Connector, interval time.Duration, handle func(context.Context, model.TelemetryFrame) error) error {
	if connector == nil || handle == nil {
		return model.ErrInvalid
	}
	if err := connector.Open(ctx); err != nil {
		return err
	}
	defer connector.Close()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			f, e := connector.Read(ctx)
			if e == io.EOF {
				return nil
			}
			if e != nil {
				return e
			}
			if e := handle(ctx, f); e != nil {
				return e
			}
		}
	}
}
