package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

func WaterBalance(tower string, makeup, blowdown, evaporation float64) model.WaterBalance {
	unaccounted := makeup - blowdown - evaporation
	return model.WaterBalance{TowerID: tower, Makeup: makeup, Blowdown: blowdown, Evaporation: evaporation, Unaccounted: unaccounted}
}
func BalanceStatus(b model.WaterBalance, tolerance float64) string {
	v := b.Unaccounted
	if v < 0 {
		v = -v
	}
	if v <= tolerance {
		return "balanced"
	}
	if b.Unaccounted > 0 {
		return "excess-makeup"
	}
	return "loss"
}

// IntegrateFlow calculates volume from m3/h readings using trapezoidal
// integration. Large collection gaps are deliberately excluded because a
// gateway outage must not be mistaken for a stable flow measurement.
func IntegrateFlow(sensor string, readings []model.Reading, maxGap time.Duration) model.FlowCoverage {
	coverage := model.FlowCoverage{Sensor: sensor, Samples: len(readings)}
	if len(readings) < 2 || maxGap <= 0 {
		return coverage
	}
	rows := append([]model.Reading(nil), readings...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].RecordedAt.Before(rows[j].RecordedAt) })
	for i := 1; i < len(rows); i++ {
		previous, current := rows[i-1], rows[i]
		interval := current.RecordedAt.Sub(previous.RecordedAt)
		if interval <= 0 {
			coverage.GapCount++
			continue
		}
		// A collection gap longer than maxGap means the gateway was offline
		// between samples; the missing span must not be integrated as if a
		// steady flow had persisted, otherwise offline periods inflate the
		// cumulative volume.
		if interval > maxGap {
			coverage.GapCount++
			continue
		}
		if previous.Value < 0 || current.Value < 0 {
			coverage.GapCount++
			continue
		}
		coverage.CoveredSecond += interval.Seconds()
		coverage.Volume += ((previous.Value + current.Value) / 2) * interval.Hours()
	}
	return coverage
}
