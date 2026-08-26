package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
)

func Evaluate(w model.WindowSummary, t model.Threshold) (string, string, string) {
	v := w.Maximum
	if t.Direction == "below" {
		v = w.Minimum
		if v <= t.Critical {
			return model.SeverityCritical, fmt.Sprintf("%s below critical %.2f", t.Sensor, t.Critical), "critical"
		}
		if v <= t.Warning {
			return model.SeverityWarning, fmt.Sprintf("%s below warning %.2f", t.Sensor, t.Warning), "warning"
		}
		return "", "", ""
	}
	if v >= t.Critical {
		return model.SeverityCritical, fmt.Sprintf("%s above critical %.2f", t.Sensor, t.Critical), "critical"
	}
	if v >= t.Warning {
		return model.SeverityWarning, fmt.Sprintf("%s above warning %.2f", t.Sensor, t.Warning), "warning"
	}
	return "", "", ""
}
func EnergyIndex(t model.Tower, w model.WindowSummary) float64 {
	if t.DesignKW <= 0 || w.Average <= 0 {
		return 0
	}
	return t.DesignKW / w.Average
}
