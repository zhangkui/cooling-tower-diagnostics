package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"math"
	"sort"
	"time"
)

type MatrixCell struct {
	Sensor      string
	Metric      string
	Value       float64
	Status      string
	Explanation string
}

type DiagnosticMatrix struct {
	TowerID string
	From    time.Time
	To      time.Time
	Cells   []MatrixCell
}

func BuildMatrix(tower string, readings []model.Reading, thresholds map[string]model.Threshold, from, to time.Time) DiagnosticMatrix {
	grouped := map[string][]model.Reading{}
	for _, r := range readings {
		if r.TowerID != tower {
			continue
		}
		if r.RecordedAt.Before(from) || r.RecordedAt.After(to) {
			continue
		}
		grouped[r.Sensor] = append(grouped[r.Sensor], r)
	}
	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	m := DiagnosticMatrix{TowerID: tower, From: from, To: to, Cells: []MatrixCell{}}
	for _, sensor := range keys {
		rows := grouped[sensor]
		var sum, min, max float64
		min, max = rows[0].Value, rows[0].Value
		for _, r := range rows {
			sum += r.Value
			if r.Value < min {
				min = r.Value
			}
			if r.Value > max {
				max = r.Value
			}
		}
		avg := sum / float64(len(rows))
		m.Cells = append(m.Cells,
			MatrixCell{Sensor: sensor, Metric: "average", Value: avg, Status: "normal", Explanation: fmt.Sprintf("%d samples", len(rows))},
			MatrixCell{Sensor: sensor, Metric: "minimum", Value: min, Status: "normal"},
			MatrixCell{Sensor: sensor, Metric: "maximum", Value: max, Status: "normal"},
			MatrixCell{Sensor: sensor, Metric: "spread", Value: max - min, Status: "normal"},
		)
		if t, ok := thresholds[sensor]; ok {
			for i := range m.Cells {
				if m.Cells[i].Sensor != sensor {
					continue
				}
				if t.Direction == "above" && m.Cells[i].Value >= t.Critical {
					m.Cells[i].Status = model.SeverityCritical
				}
				if t.Direction == "below" && m.Cells[i].Value <= t.Critical {
					m.Cells[i].Status = model.SeverityCritical
				}
			}
		}
	}
	return m
}

func MatrixRisk(m DiagnosticMatrix) string {
	risk := 0.0
	for _, c := range m.Cells {
		if c.Status == model.SeverityCritical {
			risk += 0.5
		}
		if c.Status == model.SeverityWarning {
			risk += 0.2
		}
	}
	if math.Min(risk, 1) >= 0.8 {
		return model.SeverityCritical
	}
	if risk > 0 {
		return model.SeverityWarning
	}
	return "normal"
}
