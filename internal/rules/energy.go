package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"math"
	"sort"
	"time"
)

func EnergySamples(tower string, rows []model.Reading) []model.EnergySample {
	by := map[time.Time]map[string]float64{}
	for _, r := range rows {
		if r.TowerID != tower {
			continue
		}
		bucket := r.RecordedAt.Truncate(time.Minute)
		if by[bucket] == nil {
			by[bucket] = map[string]float64{}
		}
		by[bucket][r.Sensor] = r.Value
	}
	keys := make([]time.Time, 0, len(by))
	for k := range by {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
	out := []model.EnergySample{}
	for _, k := range keys {
		v := by[k]
		in := v["input_kw"]
		outv := v["output_kw"]
		eff := 0.0
		if in > 0 {
			eff = outv / in
		}
		out = append(out, model.EnergySample{TowerID: tower, At: k, InputKW: in, OutputKW: outv, Efficiency: eff})
	}
	return out
}
func EnergyAverage(rows []model.EnergySample) float64 {
	if len(rows) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range rows {
		sum += r.Efficiency
	}
	return math.Round(sum/float64(len(rows))*10000) / 10000
}
