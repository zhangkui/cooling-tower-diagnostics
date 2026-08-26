package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"sync"
	"time"
)

type WindowStore struct {
	mu      sync.RWMutex
	windows map[string][]model.Reading
	max     int
}

func NewWindowStore(max int) *WindowStore {
	if max < 1 {
		max = 1000
	}
	return &WindowStore{windows: map[string][]model.Reading{}, max: max}
}
func (w *WindowStore) Add(r model.Reading) {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := r.TowerID + "/" + r.Sensor
	rows := append(w.windows[key], r)
	if len(rows) > w.max {
		rows = rows[len(rows)-w.max:]
	}
	w.windows[key] = rows
}
func (w *WindowStore) List(tower, sensor string) []model.Reading {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return append([]model.Reading(nil), w.windows[tower+"/"+sensor]...)
}
func (w *WindowStore) Between(tower, sensor string, from, to time.Time) []model.Reading {
	rows := w.List(tower, sensor)
	out := []model.Reading{}
	for _, r := range rows {
		if !r.RecordedAt.Before(from) && !r.RecordedAt.After(to) {
			out = append(out, r)
		}
	}
	return out
}
func (w *WindowStore) Clear(tower string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for key := range w.windows {
		if len(key) >= len(tower) && key[:len(tower)] == tower {
			delete(w.windows, key)
		}
	}
}
func (w *WindowStore) Size() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	n := 0
	for _, rows := range w.windows {
		n += len(rows)
	}
	return n
}
