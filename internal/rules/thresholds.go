package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"strings"
)

type Result struct {
	Severity string
	Breached bool
	Message  string
	Value    float64
	Limit    float64
}

func Check(value float64, t model.Threshold) Result {
	if t.Direction == "below" {
		if value <= t.Critical {
			return Result{model.SeverityCritical, true, "below critical", value, t.Critical}
		}
		if value <= t.Warning {
			return Result{model.SeverityWarning, true, "below warning", value, t.Warning}
		}
		return Result{Value: value}
	}
	if value >= t.Critical {
		return Result{model.SeverityCritical, true, "above critical", value, t.Critical}
	}
	if value >= t.Warning {
		return Result{model.SeverityWarning, true, "above warning", value, t.Warning}
	}
	return Result{Value: value}
}
func NormalizeDirection(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "lt" || v == "below" || v == "lower" {
		return "below"
	}
	return "above"
}
func ThresholdFor(sensor string, defaults map[string]model.Threshold) model.Threshold {
	if t, ok := defaults[sensor]; ok {
		return t
	}
	return model.Threshold{Sensor: sensor, Warning: 0, Critical: 0, Direction: "above", WindowMinutes: 15}
}
