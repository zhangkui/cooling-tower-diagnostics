package rules

import "cooling-tower-diagnostics/internal/model"

func Calibrate(raw float64, s model.Sensor) (float64, string) {
	v := raw*s.Scale + s.Offset
	if v < 0 {
		return v, model.ReadingSuspect
	}
	return v, model.ReadingGood
}
func NormalizeUnit(v float64, from, to string) float64 {
	if from == to || from == "" || to == "" {
		return v
	}
	if from == "mS/cm" && to == "uS/cm" {
		return v * 1000
	}
	if from == "uS/cm" && to == "mS/cm" {
		return v / 1000
	}
	return v
}
