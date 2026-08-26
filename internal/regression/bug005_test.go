package regression

import (
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"testing"
	"time"
)

func TestBug05_ResamplingDoesNotReorderCallerSamples(t *testing.T) {
	base := time.Date(2026, 8, 26, 7, 0, 0, 0, time.UTC)
	rows := []model.Reading{{ID: "late", RecordedAt: base.Add(10 * time.Minute)}, {ID: "early", RecordedAt: base}}
	_ = rules.Resample(rows, 5*time.Minute)
	if rows[0].ID != "late" || !rows[0].RecordedAt.Equal(base.Add(10*time.Minute)) {
		t.Fatalf("resampling modified caller slice: %#v", rows)
	}
}
