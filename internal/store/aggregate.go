package store

import (
	"context"
	"time"
)

type AggregateStore struct{ DB *DB }
type AggregateRow struct {
	Sensor  string
	Count   int
	Average float64
	Minimum float64
	Maximum float64
}

func (s AggregateStore) Daily(ctx context.Context, tower string, day time.Time) ([]AggregateRow, error) {
	from := day.Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT sensor,COUNT(*),AVG(value),MIN(value),MAX(value) FROM readings WHERE tower_id=? AND recorded_at>=? AND recorded_at<? GROUP BY sensor ORDER BY sensor`, tower, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano))
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []AggregateRow{}
	for rows.Next() {
		var a AggregateRow
		if e := rows.Scan(&a.Sensor, &a.Count, &a.Average, &a.Minimum, &a.Maximum); e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
func (s AggregateStore) Hourly(ctx context.Context, tower string, from, to time.Time) ([]AggregateRow, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT sensor,COUNT(*),AVG(value),MIN(value),MAX(value) FROM readings WHERE tower_id=? AND recorded_at>=? AND recorded_at<=? GROUP BY sensor ORDER BY sensor`, tower, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano))
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []AggregateRow{}
	for rows.Next() {
		var a AggregateRow
		if e := rows.Scan(&a.Sensor, &a.Count, &a.Average, &a.Minimum, &a.Maximum); e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
