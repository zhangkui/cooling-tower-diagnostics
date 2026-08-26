package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"strings"
	"time"
)

type CommandResult struct {
	ID      string
	Status  string
	Message string
	At      time.Time
}

func (s *Service) PauseTower(ctx context.Context, id string) (CommandResult, error) {
	return s.changeTower(ctx, id, model.TowerPaused, "paused")
}
func (s *Service) ResumeTower(ctx context.Context, id string) (CommandResult, error) {
	return s.changeTower(ctx, id, model.TowerOnline, "resumed")
}
func (s *Service) changeTower(ctx context.Context, id, status, event string) (CommandResult, error) {
	if err := s.SetTowerStatus(ctx, id, status); err != nil {
		return CommandResult{}, err
	}
	return CommandResult{ID: id, Status: status, Message: event, At: s.Clock.Now()}, nil
}

func (s *Service) ValidateAndIngest(ctx context.Context, req model.ReadingRequest) (model.Reading, error) {
	if strings.TrimSpace(req.Unit) == "" {
		return model.Reading{}, fmt.Errorf("unit required: %w", model.ErrInvalid)
	}
	return s.IngestReading(ctx, req)
}

func (s *Service) IngestMany(ctx context.Context, tower string, rows []model.ReadingRequest) (int, error) {
	count := 0
	for _, row := range rows {
		if row.TowerID == "" {
			row.TowerID = tower
		}
		if _, err := s.ValidateAndIngest(ctx, row); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) EvaluateAll(ctx context.Context, towerIDs []string) (map[string]model.Diagnosis, error) {
	out := map[string]model.Diagnosis{}
	for _, id := range towerIDs {
		d, err := s.EvaluateAndAlert(ctx, model.DiagnosisRequest{TowerID: id, SinceMinutes: 30})
		if err != nil {
			return out, err
		}
		out[id] = d
	}
	return out, nil
}
