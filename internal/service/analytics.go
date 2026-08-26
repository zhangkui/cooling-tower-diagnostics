package service

import (
	"context"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type Analytics struct {
	Queries  store.QueryStore
	Towers   store.TowerStore
	Readings store.ReadingStore
	Engine   diagnostics.Engine
}

func (a Analytics) Trend(ctx context.Context, tower, sensor string, minutes int) (model.Trend, error) {
	rows, e := a.Readings.Recent(ctx, tower, sensor, time.Now().UTC().Add(-time.Duration(minutes)*time.Minute))
	if e != nil {
		return model.Trend{}, e
	}
	return rules.Trend(rows, 0), nil
}
func (a Analytics) Anomalies(ctx context.Context, tower, sensor string, minutes int) ([]model.Anomaly, error) {
	rows, e := a.Readings.Recent(ctx, tower, sensor, time.Now().UTC().Add(-time.Duration(minutes)*time.Minute))
	if e != nil {
		return nil, e
	}
	return rules.DetectAnomalies(rows, 0, 10), nil
}
func (a Analytics) Quality(ctx context.Context, tower string) (map[string]int64, error) {
	return a.Queries.CountByQuality(ctx, tower)
}
func (a Analytics) Checklist(ctx context.Context, tower string) ([]diagnostics.Check, error) {
	t, e := a.Towers.Get(ctx, tower)
	if e != nil {
		return nil, e
	}
	rows := []model.Reading{}
	for _, sensor := range model.SensorKinds() {
		r, _ := a.Readings.Recent(ctx, tower, sensor, time.Now().UTC().Add(-time.Hour))
		rows = append(rows, r...)
	}
	return diagnostics.RunChecklist(ctx, rows, t), nil
}
