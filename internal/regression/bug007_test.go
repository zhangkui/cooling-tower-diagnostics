package regression

import (
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"testing"
	"time"
)

func TestBug07_LargeMeterGapIsExcludedFromVolume(t *testing.T) {
	start := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	rows := []model.Reading{{Value: 10, RecordedAt: start}, {Value: 10, RecordedAt: start.Add(time.Hour)}}
	coverage := rules.IntegrateFlow("makeup_flow", rows, 15*time.Minute)
	if coverage.Volume != 0 || coverage.GapCount != 1 { t.Fatalf("coverage=%+v", coverage) }
}
