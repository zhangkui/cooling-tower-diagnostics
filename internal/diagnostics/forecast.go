package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"time"
)

type Forecaster struct {
	Readings store.ReadingStore
	Model    rules.ForecastModel
}

func (f Forecaster) Run(ctx context.Context, tower, sensor string, horizon int) (model.Forecast, error) {
	rows, e := f.Readings.Recent(ctx, tower, sensor, time.Now().UTC().Add(-24*time.Hour))
	if e != nil {
		return model.Forecast{}, e
	}
	if err := ctx.Err(); err != nil {
		return model.Forecast{}, err
	}
	return f.Model.Predict(rows, horizon, time.Now().UTC()), nil
}
func AlertFromForecast(f model.Forecast, limit float64) bool {
	for _, v := range f.Upper {
		if v >= limit {
			return true
		}
	}
	return false
}
func ForecastRange(f model.Forecast) (float64, float64) {
	if len(f.Predicted) == 0 {
		return 0, 0
	}
	lo, hi := f.Predicted[0], f.Predicted[0]
	for _, v := range f.Predicted {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return lo, hi
}
