package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"time"
)

type Check struct {
	Code    string
	Passed  bool
	Message string
}

func RunChecklist(ctx context.Context, rows []model.Reading, tower model.Tower) []Check {
	out := []Check{}
	select {
	case <-ctx.Done():
		return []Check{{Code: "CTX", Message: ctx.Err().Error()}}
	default:
	}
	if tower.Status == model.TowerOnline {
		out = append(out, Check{Code: "STATUS", Passed: true, Message: "tower online"})
	} else {
		out = append(out, Check{Code: "STATUS", Message: "tower paused"})
	}
	if len(rows) > 0 {
		out = append(out, Check{Code: "DATA", Passed: true, Message: "telemetry present"})
	} else {
		out = append(out, Check{Code: "DATA", Message: "no telemetry"})
	}
	gaps := rules.FindGaps(rows, time.Now().Add(-time.Hour), time.Now(), 5*time.Minute)
	out = append(out, Check{Code: "GAPS", Passed: len(gaps) == 0, Message: "gap scan complete"})
	return out
}
