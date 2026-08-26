package service

import (
	"bytes"
	"context"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"fmt"
	"time"
)

type Exporter struct {
	Queries store.QueryStore
	Engine  diagnostics.Engine
}

func (e Exporter) CSV(ctx context.Context, f model.ReportFilter) ([]byte, model.ExportResult, error) {
	rows, err := e.Queries.ReadingsByRange(ctx, f)
	if err != nil {
		return nil, model.ExportResult{}, err
	}
	d := model.DiagnosticReport{ID: fmt.Sprintf("report-%d", time.Now().UnixNano()), Filter: f, Rows: make([]model.ReportRow, 0, len(rows)), CreatedAt: time.Now().UTC()}
	for _, r := range rows {
		d.Rows = append(d.Rows, model.ReportRow{Timestamp: r.RecordedAt, Sensor: r.Sensor, Value: r.Value, Quality: r.Quality})
	}
	var b bytes.Buffer
	if err := diagnostics.WriteCSV(&b, d); err != nil {
		return nil, model.ExportResult{}, err
	}
	return b.Bytes(), model.ExportResult{Format: "csv", Filename: d.ID + ".csv", Bytes: b.Len(), CreatedAt: d.CreatedAt}, nil
}
func (e Exporter) JSON(ctx context.Context, f model.ReportFilter) (model.DiagnosticReport, error) {
	r := diagnostics.Reporter{Readings: store.ReadingStore{DB: e.Queries.DB}, Engine: e.Engine}
	return r.Build(ctx, f)
}
