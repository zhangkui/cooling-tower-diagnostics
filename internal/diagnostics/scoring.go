package diagnostics

import (
	"cooling-tower-diagnostics/internal/model"
	"math"
)

func Score(d model.Diagnosis) float64 {
	base := 0.0
	switch d.Risk {
	case model.SeverityCritical:
		base = 0.9
	case model.SeverityWarning:
		base = 0.6
	default:
		base = 0.1
	}
	if d.EnergyIndex > 1.2 {
		base += 0.1
	}
	if base > 1 {
		base = 1
	}
	return math.Round(base*100) / 100
}
