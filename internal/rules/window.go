package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

func Summarize(readings []model.Reading, from, to time.Time) model.WindowSummary {
	w := model.WindowSummary{From: from, To: to}
	if len(readings) == 0 {
		return w
	}
	sort.Slice(readings, func(i, j int) bool { return readings[i].RecordedAt.Before(readings[j].RecordedAt) })
	w.TowerID = readings[0].TowerID
	w.Sensor = readings[0].Sensor
	w.Count = len(readings)
	w.Minimum = readings[0].Value
	w.Maximum = readings[0].Value
	for _, r := range readings {
		w.Average += r.Value
		if r.Value < w.Minimum {
			w.Minimum = r.Value
		}
		if r.Value > w.Maximum {
			w.Maximum = r.Value
		}
	}
	w.Average /= float64(w.Count)
	w.Spread = w.Maximum - w.Minimum
	return w
}
