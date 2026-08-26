package diagnostics

import "cooling-tower-diagnostics/internal/model"

func Explain(d model.Diagnosis) []string {
	out := append([]string{}, d.Findings...)
	if d.EnergyIndex > 1.2 {
		out = append(out, "energy efficiency is above the expected design ratio")
	}
	if len(out) == 0 {
		out = append(out, "no active rule violation in the selected window")
	}
	return out
}
