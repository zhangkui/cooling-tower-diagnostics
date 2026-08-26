package service

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"strings"
	"time"
)

func FormatReading(r model.Reading) string {
	return fmt.Sprintf("%s %s %.3f %s [%s]", r.RecordedAt.Format(time.RFC3339), r.Sensor, r.Value, r.Unit, r.Quality)
}
func FormatTower(t model.Tower) string {
	return fmt.Sprintf("%s (%s) %s design=%.1fkW", t.Name, t.ID, t.Status, t.DesignKW)
}
func FormatAlert(a model.Alert) string {
	return fmt.Sprintf("%s %s %s %s", a.ID, a.Severity, a.State, a.Message)
}
func JoinFindings(findings []string) string { return strings.Join(findings, " | ") }
func Truncate(value string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(value)
	if len(r) <= max {
		return value
	}
	return string(r[:max]) + "..."
}
