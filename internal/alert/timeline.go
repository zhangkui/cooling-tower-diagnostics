package alert

import (
	"cooling-tower-diagnostics/internal/model"
	"sort"
	"time"
)

type TimelinePoint struct {
	At       time.Time
	State    string
	Severity string
}

func Timeline(alerts []model.Alert) []TimelinePoint {
	out := []TimelinePoint{}
	for _, a := range alerts {
		out = append(out, TimelinePoint{a.OpenedAt, a.State, a.Severity})
		if a.AcknowledgedAt != nil {
			out = append(out, TimelinePoint{*a.AcknowledgedAt, model.AlertAcknowledged, a.Severity})
		}
		if a.ClosedAt != nil {
			out = append(out, TimelinePoint{*a.ClosedAt, model.AlertClosed, a.Severity})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
