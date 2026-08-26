package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

type Bucket struct {
	Start   time.Time
	End     time.Time
	Count   int
	Average float64
	Minimum float64
	Maximum float64
}

func Bucketize(readings []model.Reading, size time.Duration) []Bucket {
	if len(readings) == 0 || size <= 0 {
		return nil
	}
	copyRows := append([]model.Reading(nil), readings...)
	sort.Slice(copyRows, func(i, j int) bool { return copyRows[i].RecordedAt.Before(copyRows[j].RecordedAt) })
	out := []Bucket{}
	var cur *Bucket
	for _, r := range copyRows {
		start := r.RecordedAt.Truncate(size)
		if cur == nil || !start.Equal(cur.Start) {
			if cur != nil {
				cur.Average /= float64(cur.Count)
				out = append(out, *cur)
			}
			cur = &Bucket{Start: start, End: start.Add(size), Minimum: r.Value, Maximum: r.Value}
		}
		cur.Count++
		cur.Average += r.Value
		if r.Value < cur.Minimum {
			cur.Minimum = r.Value
		}
		if r.Value > cur.Maximum {
			cur.Maximum = r.Value
		}
	}
	if cur != nil {
		cur.Average /= float64(cur.Count)
		out = append(out, *cur)
	}
	return out
}
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	v := append([]float64(nil), values...)
	sort.Float64s(v)
	if p <= 0 {
		return v[0]
	}
	if p >= 1 {
		return v[len(v)-1]
	}
	i := p * float64(len(v)-1)
	lo := int(i)
	hi := lo + 1
	if hi >= len(v) {
		return v[lo]
	}
	return v[lo] + (v[hi]-v[lo])*(i-float64(lo))
}
