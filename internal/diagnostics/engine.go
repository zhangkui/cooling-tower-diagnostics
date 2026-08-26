package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/clock"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"fmt"
	"time"
)

type Engine struct {
	Towers     store.TowerStore
	Readings   store.ReadingStore
	Thresholds store.ThresholdStore
	Clock      clock.Clock
}

func (e Engine) Run(ctx context.Context, towerID string, minutes int) (model.Diagnosis, error) {
	if minutes <= 0 {
		minutes = 30
	}
	tower, err := e.Towers.Get(ctx, towerID)
	if err != nil {
		return model.Diagnosis{}, err
	}
	ts, err := e.Thresholds.List(ctx)
	if err != nil {
		return model.Diagnosis{}, err
	}
	d := model.Diagnosis{TowerID: towerID, GeneratedAt: e.Clock.Now(), Findings: []string{}}
	for _, t := range ts {
		rs, er := e.Readings.Recent(ctx, towerID, t.Sensor, e.Clock.Now().Add(-time.Duration(minutes)*time.Minute))
		if er != nil {
			return d, er
		}
		w := rules.Summarize(rs, e.Clock.Now().Add(-time.Duration(minutes)*time.Minute), e.Clock.Now())
		if w.Count == 0 {
			continue
		}
		if sev, msg, _ := rules.Evaluate(w, t); sev != "" {
			d.Findings = append(d.Findings, fmt.Sprintf("%s: %s", sev, msg))
			if sev == model.SeverityCritical {
				d.Risk = model.SeverityCritical
			} else if d.Risk == "" {
				d.Risk = sev
			}
		}
		if d.EnergyIndex == 0 {
			d.EnergyIndex = rules.EnergyIndex(tower, w)
		}
	}
	if d.Risk == "" {
		d.Risk = "normal"
	}
	return d, nil
}
