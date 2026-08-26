package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type CalibrationStore struct{ DB *DB }

func (s CalibrationStore) Put(ctx context.Context, p model.CalibrationProfile) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO calibrations(sensor,version,effective_at,offset,scale,min_value,max_value) VALUES(?,?,?,?,?,?,?) ON CONFLICT(sensor) DO UPDATE SET version=excluded.version,effective_at=excluded.effective_at,offset=excluded.offset,scale=excluded.scale,min_value=excluded.min_value,max_value=excluded.max_value`, p.Sensor, p.Version, p.EffectiveAt.Format(time.RFC3339Nano), p.Offset, p.Scale, p.Min, p.Max)
	return e
}
func (s CalibrationStore) Get(ctx context.Context, sensor string) (model.CalibrationProfile, error) {
	var p model.CalibrationProfile
	var ts string
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT sensor,version,effective_at,offset,scale,min_value,max_value FROM calibrations WHERE sensor=?`, sensor).Scan(&p.Sensor, &p.Version, &ts, &p.Offset, &p.Scale, &p.Min, &p.Max)
	if e != nil {
		return p, model.ErrNotFound
	}
	p.EffectiveAt, _ = time.Parse(time.RFC3339Nano, ts)
	return p, nil
}
func (s CalibrationStore) List(ctx context.Context) ([]model.CalibrationProfile, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT sensor,version,effective_at,offset,scale,min_value,max_value FROM calibrations ORDER BY sensor`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.CalibrationProfile{}
	for rows.Next() {
		var p model.CalibrationProfile
		var ts string
		if e := rows.Scan(&p.Sensor, &p.Version, &ts, &p.Offset, &p.Scale, &p.Min, &p.Max); e != nil {
			return nil, e
		}
		p.EffectiveAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, p)
	}
	return out, rows.Err()
}
