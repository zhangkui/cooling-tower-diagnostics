package ingest

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"math"
)

type QualityChecker struct {
	Catalog map[string]model.SensorCatalog
	MaxJump float64
}

func (q QualityChecker) Check(r model.Reading, previous *model.Reading) string {
	if math.IsNaN(r.Value) || math.IsInf(r.Value, 0) {
		return model.ReadingSuspect
	}
	if cat, ok := q.Catalog[r.Sensor]; ok && (r.Value < cat.MinValue || r.Value > cat.MaxValue) {
		return model.ReadingSuspect
	}
	if previous != nil && q.MaxJump > 0 && math.Abs(r.Value-previous.Value) > q.MaxJump {
		return model.ReadingSuspect
	}
	return model.ReadingGood
}

func (q QualityChecker) Explain(r model.Reading) string {
	if r.Quality == model.ReadingGood {
		return "within configured quality bounds"
	}
	return fmt.Sprintf("suspect %s reading %.3f", r.Sensor, r.Value)
}
