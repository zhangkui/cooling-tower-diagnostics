package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
)

type ThresholdStore struct{ DB *DB }

func (s ThresholdStore) Put(ctx context.Context, t model.Threshold) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO thresholds(sensor,warning,critical,direction,window_minutes) VALUES(?,?,?,?,?) ON CONFLICT(sensor) DO UPDATE SET warning=excluded.warning,critical=excluded.critical,direction=excluded.direction,window_minutes=excluded.window_minutes`, t.Sensor, t.Warning, t.Critical, t.Direction, t.WindowMinutes)
	return e
}
func (s ThresholdStore) Get(ctx context.Context, sensor string) (model.Threshold, error) {
	var t model.Threshold
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT sensor,warning,critical,direction,window_minutes FROM thresholds WHERE sensor=?`, sensor).Scan(&t.Sensor, &t.Warning, &t.Critical, &t.Direction, &t.WindowMinutes)
	if e != nil {
		return t, model.ErrNotFound
	}
	return t, nil
}
func (s ThresholdStore) List(ctx context.Context) ([]model.Threshold, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT sensor,warning,critical,direction,window_minutes FROM thresholds ORDER BY sensor`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Threshold{}
	for rows.Next() {
		var t model.Threshold
		if e := rows.Scan(&t.Sensor, &t.Warning, &t.Critical, &t.Direction, &t.WindowMinutes); e != nil {
			return nil, e
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
