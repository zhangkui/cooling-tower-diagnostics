package service

import (
	"context"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"fmt"
	"time"
)

type DiagnosisService struct {
	Engine   diagnostics.Engine
	Readings store.ReadingStore
	Towers   store.TowerStore
	Metrics  *Metrics
}

func (d DiagnosisService) Run(ctx context.Context, tower string, minutes int) (model.Diagnosis, error) {
	if d.Metrics != nil {
		d.Metrics.Inc("diagnosis_runs")
	}
	start := time.Now()
	result, e := d.Engine.Run(ctx, tower, minutes)
	if d.Metrics != nil {
		d.Metrics.Observe("diagnosis_duration", time.Since(start))
	}
	if e != nil {
		return result, e
	}
	result.Findings = diagnostics.Explain(result)
	return result, nil
}
func (d DiagnosisService) Compare(ctx context.Context, tower, sensor string, minutes int) ([]diagnostics.Comparison, error) {
	now := time.Now().UTC()
	current, err := d.Readings.Recent(ctx, tower, sensor, now.Add(-time.Duration(minutes)*time.Minute))
	if err != nil {
		return nil, err
	}
	previous, err := d.Readings.Recent(ctx, tower, sensor, now.Add(-2*time.Duration(minutes)*time.Minute))
	if err != nil {
		return nil, err
	}
	return diagnostics.CompareWindows(current, previous), nil
}
func (d DiagnosisService) Matrix(ctx context.Context, tower string, minutes int) (rules.DiagnosticMatrix, error) {
	now := time.Now().UTC()
	rows := []model.Reading{}
	for _, sensor := range model.SensorKinds() {
		r, e := d.Readings.Recent(ctx, tower, sensor, now.Add(-time.Duration(minutes)*time.Minute))
		if e != nil {
			return rules.DiagnosticMatrix{}, e
		}
		rows = append(rows, r...)
	}
	thresholdRows, e := store.ThresholdStore{DB: d.Readings.DB}.List(ctx)
	if e != nil {
		return rules.DiagnosticMatrix{}, e
	}
	thresholds := map[string]model.Threshold{}
	for _, row := range thresholdRows {
		thresholds[row.Sensor] = row
	}
	return rules.BuildMatrix(tower, rows, thresholds, now.Add(-time.Duration(minutes)*time.Minute), now), nil
}
func (d DiagnosisService) ExplainError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("diagnosis failed: %v", err)
}
