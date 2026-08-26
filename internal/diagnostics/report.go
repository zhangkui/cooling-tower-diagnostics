package diagnostics

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Reporter struct {
	Readings store.ReadingStore
	Engine   Engine
}

func (r Reporter) Build(ctx context.Context, f model.ReportFilter) (model.DiagnosticReport, error) {
	d, e := r.Engine.Run(ctx, f.TowerID, int(time.Since(f.From).Minutes()))
	if e != nil {
		return model.DiagnosticReport{}, e
	}
	rows := []model.ReportRow{}
	for _, sensor := range f.Sensors {
		rs, er := r.Readings.Recent(ctx, f.TowerID, sensor, f.From)
		if er != nil {
			return model.DiagnosticReport{}, er
		}
		for _, x := range rs {
			if !f.To.IsZero() && x.RecordedAt.After(f.To) {
				continue
			}
			rows = append(rows, model.ReportRow{Timestamp: x.RecordedAt, Sensor: x.Sensor, Value: x.Value, Quality: x.Quality})
		}
	}
	return model.DiagnosticReport{ID: fmt.Sprintf("report-%d", time.Now().UnixNano()), Filter: f, Rows: rows, Findings: d.Findings, CreatedAt: time.Now().UTC()}, nil
}
func WriteCSV(w io.Writer, r model.DiagnosticReport) error {
	c := csv.NewWriter(w)
	if e := c.Write([]string{"timestamp", "sensor", "value", "quality", "alert_state"}); e != nil {
		return e
	}
	for _, x := range r.Rows {
		if e := c.Write([]string{x.Timestamp.Format(time.RFC3339), x.Sensor, strconv.FormatFloat(x.Value, 'f', 3, 64), x.Quality, x.AlertState}); e != nil {
			return e
		}
	}
	c.Flush()
	return c.Error()
}
