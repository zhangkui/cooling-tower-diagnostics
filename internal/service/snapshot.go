package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"time"
)

type Snapshot struct {
	GeneratedAt time.Time
	Towers      []model.Tower
	Alerts      []model.Alert
	Thresholds  []model.Threshold
	Metrics     map[string]any
}

func (s *Service) Snapshot(ctx context.Context) (Snapshot, error) {
	t, e := s.ListTowers(ctx)
	if e != nil {
		return Snapshot{}, e
	}
	a, e := s.ListAlerts(ctx, "")
	if e != nil {
		return Snapshot{}, e
	}
	th, e := s.ListThresholds(ctx)
	if e != nil {
		return Snapshot{}, e
	}
	return Snapshot{GeneratedAt: s.Clock.Now(), Towers: t, Alerts: a, Thresholds: th}, nil
}
func (s Snapshot) JSON() ([]byte, error) { return json.MarshalIndent(s, "", "  ") }
func (s Snapshot) Healthy() bool {
	for _, a := range s.Alerts {
		if a.State != model.AlertClosed && a.Severity == model.SeverityCritical {
			return false
		}
	}
	return true
}
