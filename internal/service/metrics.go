package service

import (
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.RWMutex
	counts    map[string]int64
	durations map[string]time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{counts: map[string]int64{}, durations: map[string]time.Duration{}}
}
func (m *Metrics) Inc(name string) { m.mu.Lock(); defer m.mu.Unlock(); m.counts[name]++ }
func (m *Metrics) Observe(name string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.durations[name] += d
}
func (m *Metrics) Snapshot() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := map[string]any{}
	for k, v := range m.counts {
		out[k] = v
	}
	for k, v := range m.durations {
		out[k+"_duration_ms"] = v.Milliseconds()
	}
	return out
}
