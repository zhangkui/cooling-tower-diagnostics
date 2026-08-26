package service

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"math"
	"strings"
	"time"
)

type Validator struct{ Clock func() time.Time }

func (v Validator) Tower(request model.CreateTowerRequest) error {
	if strings.TrimSpace(request.Name) == "" {
		return fmt.Errorf("tower name: %w", model.ErrInvalid)
	}
	if strings.TrimSpace(request.Site) == "" {
		return fmt.Errorf("tower site: %w", model.ErrInvalid)
	}
	if request.DesignKW <= 0 || math.IsNaN(request.DesignKW) {
		return fmt.Errorf("tower design power: %w", model.ErrInvalid)
	}
	return nil
}

func (v Validator) Sensor(request model.CreateSensorRequest) error {
	if request.TowerID == "" || request.Kind == "" {
		return model.ErrInvalid
	}
	if request.Unit == "" {
		return fmt.Errorf("sensor unit: %w", model.ErrInvalid)
	}
	if request.Scale == 0 {
		return fmt.Errorf("sensor scale: %w", model.ErrInvalid)
	}
	return nil
}

func (v Validator) Reading(request model.ReadingRequest) error {
	if request.TowerID == "" {
		return fmt.Errorf("tower id: %w", model.ErrInvalid)
	}
	if strings.TrimSpace(request.Sensor) == "" {
		return fmt.Errorf("sensor: %w", model.ErrInvalid)
	}
	if math.IsNaN(request.Value) || math.IsInf(request.Value, 0) {
		return fmt.Errorf("value: %w", model.ErrInvalid)
	}
	return nil
}

func (v Validator) Diagnosis(request model.DiagnosisRequest) error {
	if request.TowerID == "" {
		return model.ErrInvalid
	}
	if request.SinceMinutes < 1 || request.SinceMinutes > 24*60 {
		return model.ErrInvalid
	}
	return nil
}

func (v Validator) Report(filter model.ReportFilter) error {
	if filter.TowerID == "" {
		return model.ErrInvalid
	}
	if filter.From.IsZero() || filter.To.IsZero() {
		return model.ErrInvalid
	}
	if filter.To.Before(filter.From) {
		return fmt.Errorf("time range: %w", model.ErrInvalid)
	}
	return nil
}

func (v Validator) Timestamp(at time.Time) time.Time {
	if at.IsZero() {
		if v.Clock != nil {
			return v.Clock()
		}
		return time.Now().UTC()
	}
	return at.UTC()
}
