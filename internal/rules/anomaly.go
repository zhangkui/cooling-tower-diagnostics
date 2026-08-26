package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"math"
	"time"
)

func DetectAnomalies(rows []model.Reading, baseline, threshold float64) []model.Anomaly {
	out := []model.Anomaly{}
	for _, r := range rows {
		score := math.Abs(r.Value - baseline)
		if score >= threshold {
			out = append(out, model.Anomaly{ID: r.ID, TowerID: r.TowerID, Sensor: r.Sensor, At: r.RecordedAt, Value: r.Value, Expected: baseline, Score: score, Reason: "deviation from operating baseline"})
		}
	}
	return out
}
func ResolveBefore(items []model.Anomaly, at time.Time) []model.Anomaly {
	out := make([]model.Anomaly, len(items))
	copy(out, items)
	for i := range out {
		if out[i].At.Before(at) {
			out[i].Resolved = true
		}
	}
	return out
}
