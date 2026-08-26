package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type Calibrator struct{ Store store.CalibrationStore }

func (c Calibrator) Apply(ctx context.Context, sensor string, raw float64) (float64, string, error) {
	p, e := c.Store.Get(ctx, sensor)
	if e != nil {
		if e == model.ErrNotFound {
			return raw, model.ReadingSuspect, nil
		}
		return 0, "", e
	}
	v := raw*p.Scale + p.Offset
	if v < p.Min || v > p.Max {
		return v, model.ReadingSuspect, nil
	}
	return v, model.ReadingGood, nil
}
func (c Calibrator) Save(ctx context.Context, p model.CalibrationProfile) error {
	if p.EffectiveAt.IsZero() {
		p.EffectiveAt = time.Now().UTC()
	}
	if p.Scale == 0 {
		p.Scale = 1
	}
	return c.Store.Put(ctx, p)
}
func (c Calibrator) Convert(raw float64, p model.CalibrationProfile) float64 {
	v, _ := rules.Calibrate(raw, model.Sensor{Scale: p.Scale, Offset: p.Offset})
	return v
}
