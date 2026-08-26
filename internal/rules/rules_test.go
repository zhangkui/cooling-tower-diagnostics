package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"testing"
	"time"
)

func TestSummarize(t *testing.T) {
	now := time.Now()
	w := Summarize([]model.Reading{{TowerID: "t", Sensor: "x", Value: 2, RecordedAt: now}, {TowerID: "t", Sensor: "x", Value: 4, RecordedAt: now.Add(time.Second)}}, now.Add(-time.Minute), now)
	if w.Average != 3 || w.Spread != 2 {
		t.Fatalf("unexpected %#v", w)
	}
}
func TestEvaluate(t *testing.T) {
	sev, _, _ := Evaluate(model.WindowSummary{Maximum: 10}, model.Threshold{Sensor: "x", Warning: 5, Critical: 9, Direction: "above"})
	if sev != model.SeverityCritical {
		t.Fatal(sev)
	}
}

func TestIntegrateFlowSkipsLongGaps(t *testing.T) {
	start := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	coverage := IntegrateFlow("makeup_flow", []model.Reading{
		{Sensor: "makeup_flow", Value: 60, RecordedAt: start},
		{Sensor: "makeup_flow", Value: 120, RecordedAt: start.Add(10 * time.Minute)},
		{Sensor: "makeup_flow", Value: 120, RecordedAt: start.Add(40 * time.Minute)},
	}, 15*time.Minute)
	if coverage.GapCount != 1 || coverage.CoveredSecond != 600 || coverage.Volume != 15 {
		t.Fatalf("unexpected coverage %#v", coverage)
	}
}
