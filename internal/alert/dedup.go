package alert

import (
	"cooling-tower-diagnostics/internal/model"
	"strings"
)

func Key(tower string, t model.Threshold) string {
	return strings.Join([]string{tower, t.Sensor, t.Direction}, "/")
}
func ShouldSuppress(existing []model.Alert, rule string) bool {
	for _, a := range existing {
		if a.Rule == rule && a.State != model.AlertClosed {
			return true
		}
	}
	return false
}
