package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"math"
)

func ValidateReading(r model.Reading) error {
	if r.TowerID == "" || r.Sensor == "" {
		return model.ErrInvalid
	}
	if math.IsNaN(r.Value) || math.IsInf(r.Value, 0) {
		return fmt.Errorf("invalid numeric value: %w", model.ErrInvalid)
	}
	if r.Unit == "" {
		return fmt.Errorf("missing unit: %w", model.ErrInvalid)
	}
	return nil
}
func ValidateThreshold(t model.Threshold) error {
	if t.Sensor == "" || t.WindowMinutes < 1 {
		return model.ErrInvalid
	}
	if t.Direction != "above" && t.Direction != "below" {
		return fmt.Errorf("direction: %w", model.ErrInvalid)
	}
	if t.Direction == "above" && t.Critical < t.Warning {
		return fmt.Errorf("critical below warning: %w", model.ErrInvalid)
	}
	if t.Direction == "below" && t.Critical > t.Warning {
		return fmt.Errorf("critical above warning: %w", model.ErrInvalid)
	}
	return nil
}
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
