package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"sort"
	"strings"
)

type ThresholdManager struct {
	Store interface {
		Put(context.Context, model.Threshold) error
		List(context.Context) ([]model.Threshold, error)
	}
}

func (m ThresholdManager) Save(ctx context.Context, req model.ThresholdRequest) (model.Threshold, error) {
	t := model.Threshold{Sensor: strings.ToLower(strings.TrimSpace(req.Sensor)), Warning: req.Warning, Critical: req.Critical, Direction: req.Direction, WindowMinutes: req.WindowMinutes}
	if t.Sensor == "" || t.WindowMinutes < 1 {
		return t, model.ErrInvalid
	}
	if t.Direction != "above" && t.Direction != "below" {
		return t, fmt.Errorf("unsupported direction: %w", model.ErrInvalid)
	}
	if t.Direction == "above" && t.Warning > t.Critical {
		return t, model.ErrInvalid
	}
	if t.Direction == "below" && t.Warning < t.Critical {
		return t, model.ErrInvalid
	}
	if err := m.Store.Put(ctx, t); err != nil {
		return t, err
	}
	return t, nil
}

func (m ThresholdManager) Sorted(ctx context.Context) ([]model.Threshold, error) {
	items, err := m.Store.List(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Sensor < items[j].Sensor })
	return items, nil
}

func (m ThresholdManager) GroupByDirection(ctx context.Context) (map[string][]model.Threshold, error) {
	items, err := m.Sorted(ctx)
	if err != nil {
		return nil, err
	}
	out := map[string][]model.Threshold{}
	for _, item := range items {
		out[item.Direction] = append(out[item.Direction], item)
	}
	return out, nil
}
