package diagnostics

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"sort"
	"strings"
)

type Summary struct {
	Risk            string
	Headline        string
	Counts          map[string]int
	Recommendations []Recommendation
}

func Summarize(d model.Diagnosis) Summary {
	counts := map[string]int{}
	for _, f := range d.Findings {
		key := "info"
		if strings.Contains(strings.ToLower(f), "critical") {
			key = "critical"
		}
		if strings.Contains(strings.ToLower(f), "warning") {
			key = "warning"
		}
		counts[key]++
	}
	recs := Recommend(d)
	headline := "Operating window is stable"
	if d.Risk == model.SeverityCritical {
		headline = "Immediate operator attention required"
	}
	if d.Risk == model.SeverityWarning {
		headline = "Review developing operating deviation"
	}
	return Summary{Risk: d.Risk, Headline: headline, Counts: counts, Recommendations: recs}
}
func FormatSummary(s Summary) string {
	keys := make([]string, 0, len(s.Counts))
	for k := range s.Counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{s.Headline}
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, s.Counts[k]))
	}
	return strings.Join(parts, "; ")
}
