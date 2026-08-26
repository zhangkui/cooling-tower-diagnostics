package rules

import (
	"cooling-tower-diagnostics/internal/model"
	"math"
	"time"
)

type ForecastModel struct {
	Window      int
	Decay       float64
	TrendWeight float64
}

func (m ForecastModel) Predict(rows []model.Reading, horizon int, now time.Time) model.Forecast {
	if horizon < 1 {
		horizon = 1
	}
	if m.Window < 1 {
		m.Window = 5
	}
	if m.Decay <= 0 || m.Decay >= 1 {
		m.Decay = .8
	}
	if len(rows) == 0 {
		return model.Forecast{HorizonMinutes: horizon, GeneratedAt: now}
	}
	start := len(rows) - m.Window
	if start < 0 {
		start = 0
	}
	window := rows[start:]
	baseline := 0.0
	for _, r := range window {
		baseline += r.Value
	}
	baseline /= float64(len(window))
	slope := 0.0
	if len(window) > 1 {
		slope = (window[len(window)-1].Value - window[0].Value) / float64(len(window)-1)
	}
	out := model.Forecast{Sensor: rows[0].Sensor, HorizonMinutes: horizon, GeneratedAt: now, Predicted: make([]float64, 0, horizon), Lower: make([]float64, 0, horizon), Upper: make([]float64, 0, horizon)}
	for i := 1; i <= horizon; i++ {
		value := baseline + slope*float64(i)*m.TrendWeight
		spread := math.Abs(slope)*float64(i) + 1
		out.Predicted = append(out.Predicted, value)
		out.Lower = append(out.Lower, value-spread)
		out.Upper = append(out.Upper, value+spread)
	}
	return out
}
func Smooth(values []float64, alpha float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	if alpha <= 0 || alpha > 1 {
		alpha = .5
	}
	out := make([]float64, len(values))
	out[0] = values[0]
	for i := 1; i < len(values); i++ {
		out[i] = alpha*values[i] + (1-alpha)*out[i-1]
	}
	return out
}
func ForecastAt(f model.Forecast, index int) float64 {
	if index < 0 || index >= len(f.Predicted) {
		return 0
	}
	return f.Predicted[index]
}
func ForecastAge(f model.Forecast, now time.Time) time.Duration { return now.Sub(f.GeneratedAt) }
