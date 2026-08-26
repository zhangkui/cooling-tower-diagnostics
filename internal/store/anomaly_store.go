package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type AnomalyStore struct{ DB *DB }

func (s AnomalyStore) Put(ctx context.Context, a model.Anomaly) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO anomalies(id,tower_id,sensor,at,value,expected,score,reason,resolved) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET resolved=excluded.resolved,reason=excluded.reason`, a.ID, a.TowerID, a.Sensor, a.At.Format(time.RFC3339Nano), a.Value, a.Expected, a.Score, a.Reason, a.Resolved)
	return e
}
func (s AnomalyStore) List(ctx context.Context, tower string, onlyOpen bool) ([]model.Anomaly, error) {
	q := `SELECT id,tower_id,sensor,at,value,expected,score,reason,resolved FROM anomalies WHERE tower_id=? ORDER BY at DESC`
	args := []any{tower}
	if onlyOpen {
		q = `SELECT id,tower_id,sensor,at,value,expected,score,reason,resolved FROM anomalies WHERE tower_id=? AND resolved=0 ORDER BY at DESC`
	}
	rows, e := s.DB.SQL.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Anomaly{}
	for rows.Next() {
		var a model.Anomaly
		var ts string
		var resolved int
		if e := rows.Scan(&a.ID, &a.TowerID, &a.Sensor, &ts, &a.Value, &a.Expected, &a.Score, &a.Reason, &resolved); e != nil {
			return nil, e
		}
		a.At, _ = time.Parse(time.RFC3339Nano, ts)
		a.Resolved = resolved != 0
		out = append(out, a)
	}
	return out, rows.Err()
}
func (s AnomalyStore) Resolve(ctx context.Context, id string) error {
	_, e := s.DB.SQL.ExecContext(ctx, `UPDATE anomalies SET resolved=1 WHERE id=?`, id)
	return e
}
