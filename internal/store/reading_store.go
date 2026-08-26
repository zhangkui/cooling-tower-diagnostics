package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"time"
)

type ReadingStore struct{ DB *DB }

func (s ReadingStore) Add(ctx context.Context, r model.Reading) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO readings(id,tower_id,sensor,value,raw_value,unit,quality,recorded_at) VALUES(?,?,?,?,?,?,?,?)`, r.ID, r.TowerID, r.Sensor, r.Value, r.RawValue, r.Unit, r.Quality, r.RecordedAt.Format(time.RFC3339Nano))
	return e
}
func (s ReadingStore) Recent(ctx context.Context, tower, sensor string, since time.Time) ([]model.Reading, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT id,tower_id,sensor,value,raw_value,unit,quality,recorded_at FROM readings WHERE tower_id=? AND sensor=? AND recorded_at>=? ORDER BY recorded_at`, tower, sensor, since.Format(time.RFC3339Nano))
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Reading{}
	for rows.Next() {
		var r model.Reading
		var ts string
		if e := rows.Scan(&r.ID, &r.TowerID, &r.Sensor, &r.Value, &r.RawValue, &r.Unit, &r.Quality, &ts); e != nil {
			return nil, e
		}
		r.RecordedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, r)
	}
	return out, rows.Err()
}
func (s ReadingStore) Count(ctx context.Context, tower string) (int, error) {
	var n int
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE tower_id=?`, tower).Scan(&n)
	return n, e
}
func (s ReadingStore) Last(ctx context.Context, tower string) (model.Reading, error) {
	var r model.Reading
	var ts string
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT id,tower_id,sensor,value,raw_value,unit,quality,recorded_at FROM readings WHERE tower_id=? ORDER BY recorded_at DESC LIMIT 1`, tower).Scan(&r.ID, &r.TowerID, &r.Sensor, &r.Value, &r.RawValue, &r.Unit, &r.Quality, &ts)
	if e != nil {
		return r, fmt.Errorf("last reading: %w", e)
	}
	r.RecordedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return r, nil
}
