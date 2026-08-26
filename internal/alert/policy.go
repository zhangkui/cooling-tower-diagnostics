package alert

import (
	"cooling-tower-diagnostics/internal/model"
	"strings"
	"time"
)

type Policy struct {
	MinimumSeverity string
	QuietPeriod     time.Duration
	EscalateAfter   time.Duration
	RequiredTags    []string
}

func SeverityRank(s string) int {
	switch strings.ToLower(s) {
	case model.SeverityCritical:
		return 3
	case model.SeverityWarning:
		return 2
	case model.SeverityInfo:
		return 1
	}
	return 0
}
func (p Policy) Accept(a model.Alert) bool {
	return SeverityRank(a.Severity) >= SeverityRank(p.MinimumSeverity)
}
func (p Policy) CanReopen(a model.Alert, now time.Time) bool {
	if a.ClosedAt == nil {
		return true
	}
	return now.Sub(*a.ClosedAt) >= p.QuietPeriod
}
func MergeMessages(alerts []model.Alert) string {
	parts := make([]string, 0, len(alerts))
	for _, a := range alerts {
		parts = append(parts, a.Message)
	}
	return strings.Join(parts, "; ")
}
