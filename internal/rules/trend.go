package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"math"
	"sort"
	"time"
)

func Trend(rows []model.Reading, baseline float64) model.Trend {
	if len(rows) == 0 {
		return model.Trend{}
	}
	copyRows := append([]model.Reading(nil), rows...)
	sort.Slice(copyRows, func(i, j int) bool { return copyRows[i].RecordedAt.Before(copyRows[j].RecordedAt) })
	first := copyRows[0]
	last := copyRows[len(copyRows)-1]
	duration := last.RecordedAt.Sub(first.RecordedAt).Hours()
	slope := 0.0
	if duration > 0 {
		slope = (last.Value - first.Value) / duration
	}
	direction := "flat"
	if slope > 0.01 {
		direction = "up"
	}
	if slope < -0.01 {
		direction = "down"
	}
	points := []model.TrendPoint{}
	for _, r := range copyRows {
		points = append(points, model.TrendPoint{At: r.RecordedAt, Value: r.Value, Baseline: baseline, Delta: r.Value - baseline})
	}
	confidence := 1 - math.Min(1, math.Abs(slope)/(math.Abs(baseline)+1))
	return model.Trend{Sensor: first.Sensor, Direction: direction, Slope: slope, Confidence: confidence, Points: points}
}
func Resample(rows []model.Reading, step time.Duration) []model.Reading {
	if step <= 0 {
		return rows
	}
	if len(rows) == 0 {
		return nil
	}
	rows = append([]model.Reading(nil), rows...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].RecordedAt.Before(rows[j].RecordedAt) })
	out := []model.Reading{}
	next := rows[0].RecordedAt
	for _, r := range rows {
		for next.Before(r.RecordedAt) {
			next = next.Add(step)
		}
		r.RecordedAt = next
		out = append(out, r)
	}
	return out
}
