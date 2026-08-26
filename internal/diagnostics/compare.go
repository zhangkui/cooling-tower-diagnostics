package diagnostics

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

type Comparison struct {
	Sensor    string
	Current   float64
	Previous  float64
	Delta     float64
	Direction string
}

func CompareWindows(current, previous []model.Reading) []Comparison {
	cm := map[string][]float64{}
	pm := map[string][]float64{}
	for _, r := range current {
		cm[r.Sensor] = append(cm[r.Sensor], r.Value)
	}
	for _, r := range previous {
		pm[r.Sensor] = append(pm[r.Sensor], r.Value)
	}
	keys := []string{}
	for k := range cm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []Comparison{}
	for _, k := range keys {
		c := avg(cm[k])
		p := avg(pm[k])
		dir := "flat"
		if c > p {
			dir = "up"
		}
		if c < p {
			dir = "down"
		}
		out = append(out, Comparison{Sensor: k, Current: c, Previous: p, Delta: c - p, Direction: dir})
	}
	return out
}
func avg(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}
func Period(now time.Time, minutes int) (time.Time, time.Time) {
	if minutes < 1 {
		minutes = 30
	}
	return now.Add(-time.Duration(minutes) * time.Minute), now
}
