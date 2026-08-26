package diagnostics

import (
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func FormatDiagnosis(d model.Diagnosis) string {
	parts := []string{fmt.Sprintf("tower=%s", d.TowerID), fmt.Sprintf("risk=%s", d.Risk), fmt.Sprintf("energy=%.3f", d.EnergyIndex)}
	parts = append(parts, d.Findings...)
	return strings.Join(parts, "; ")
}
func DiagnosisJSON(d model.Diagnosis) ([]byte, error) { return json.MarshalIndent(d, "", "  ") }
func WindowLabel(from, to time.Time) string {
	return from.Format("2006-01-02 15:04") + " - " + to.Format("15:04")
}
func RiskColor(risk string) string {
	switch risk {
	case model.SeverityCritical:
		return "#b42318"
	case model.SeverityWarning:
		return "#b54708"
	default:
		return "#067647"
	}
}
