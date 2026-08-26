package alert

import (
	"cooling-tower-diagnostics/internal/model"
	"testing"
	"time"
)

func TestLifecycle(t *testing.T) {
	n := time.Now()
	a := Open(n, "t", "r", model.SeverityWarning, "x")
	if Acknowledge(&a, n.Add(time.Minute)) != nil {
		t.Fatal("ack")
	}
}
