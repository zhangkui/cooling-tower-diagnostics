package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"fmt"
	"time"
)

type QueryStore struct{ DB *DB }

func (s QueryStore) ReadingsByRange(ctx context.Context, f model.ReportFilter) ([]model.Reading, error) {
	q := `SELECT id,tower_id,sensor,value,raw_value,unit,quality,recorded_at FROM readings WHERE tower_id=? AND recorded_at>? AND recorded_at<=? ORDER BY recorded_at`
	rows, err := s.DB.SQL.QueryContext(ctx, q, f.TowerID, f.From.Format(time.RFC3339Nano), f.To.Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Reading{}
	for rows.Next() {
		var r model.Reading
		var ts string
		if err := rows.Scan(&r.ID, &r.TowerID, &r.Sensor, &r.Value, &r.RawValue, &r.Unit, &r.Quality, &ts); err != nil {
			return nil, err
		}
		r.RecordedAt, _ = time.Parse(time.RFC3339Nano, ts)
		if len(f.Sensors) > 0 && !contains(f.Sensors, r.Sensor) {
			continue
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func contains(values []string, needle string) bool {
	for _, v := range values {
		if v == needle {
			return true
		}
	}
	return false
}

func (s QueryStore) CountByQuality(ctx context.Context, tower string) (map[string]int64, error) {
	rows, err := s.DB.SQL.QueryContext(ctx, `SELECT quality,COUNT(*) FROM readings WHERE tower_id=? GROUP BY quality`, tower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var key string
		var count int64
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		out[key] = count
	}
	return out, rows.Err()
}

func (s QueryStore) LatestForSensors(ctx context.Context, tower string, sensors []string) ([]model.Reading, error) {
	out := []model.Reading{}
	for _, sensor := range sensors {
		r, err := ReadingStore{DB: s.DB}.Recent(ctx, tower, sensor, time.Now().UTC().Add(-365*24*time.Hour))
		if err != nil {
			return nil, err
		}
		if len(r) > 0 {
			out = append(out, r[len(r)-1])
		}
	}
	return out, nil
}

func scanOptionalTime(value sql.NullString) *time.Time {
	if !value.Valid {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, value.String)
	if err != nil {
		return nil
	}
	return &t
}
func describeQueryError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("query store: %w", err)
}
