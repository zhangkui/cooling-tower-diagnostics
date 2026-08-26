package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

func FindGaps(rows []model.Reading, from, to time.Time, interval time.Duration) []model.DataGap {
	if interval <= 0 {
		return nil
	}
	group := map[string][]model.Reading{}
	for _, r := range rows {
		if r.RecordedAt.Before(from) || r.RecordedAt.After(to) {
			continue
		}
		group[r.Sensor] = append(group[r.Sensor], r)
	}
	out := []model.DataGap{}
	expected := int(to.Sub(from)/interval) + 1
	for sensor, rs := range group {
		sort.Slice(rs, func(i, j int) bool { return rs[i].RecordedAt.Before(rs[j].RecordedAt) })
		actual := len(rs)
		missing := expected - actual
		if missing < 0 {
			missing = 0
		}
		if missing > 0 {
			out = append(out, model.DataGap{TowerID: rs[0].TowerID, Sensor: sensor, From: from, To: to, Expected: expected, Actual: actual, Missing: missing})
		}
	}
	return out
}
func Interpolate(a, b model.Reading, at time.Time) model.Reading {
	if b.RecordedAt.Equal(a.RecordedAt) {
		a.RecordedAt = at
		return a
	}
	ratio := float64(at.Sub(a.RecordedAt)) / float64(b.RecordedAt.Sub(a.RecordedAt))
	a.ID = ""
	a.RecordedAt = at
	a.Value = a.Value + (b.Value-a.Value)*ratio
	a.RawValue = a.Value
	return a
}
