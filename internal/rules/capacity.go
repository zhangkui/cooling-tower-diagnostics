package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"math"
)

func CapacityStatus(tower string, rated, current float64) model.Capacity {
	util := 0.0
	if rated > 0 {
		util = current / rated
	}
	head := rated - current
	status := "normal"
	if rated <= 0 {
		status = "unconfigured"
	}
	if util >= 1 {
		status = "overload"
	} else if util >= .85 {
		status = "near-limit"
	}
	return model.Capacity{TowerID: tower, RatedFlow: rated, CurrentFlow: current, Utilization: util, Headroom: head, Status: status}
}
func CapacityMessage(c model.Capacity) string {
	switch c.Status {
	case "overload":
		return fmt.Sprintf("flow overload %.1f%%", c.Utilization*100)
	case "near-limit":
		return fmt.Sprintf("flow near limit %.1f%%", c.Utilization*100)
	case "unconfigured":
		return "rated flow is not configured"
	default:
		return fmt.Sprintf("flow headroom %.2f", c.Headroom)
	}
}
func CapacityScore(c model.Capacity) float64 {
	if c.RatedFlow <= 0 {
		return 0
	}
	return math.Max(0, 1-c.Utilization)
}
