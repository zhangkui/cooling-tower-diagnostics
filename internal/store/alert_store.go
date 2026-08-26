package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"errors"
	"time"
)

type AlertStore struct{ DB *DB }

func (s AlertStore) Upsert(ctx context.Context, a model.Alert) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO alerts(id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at) VALUES(?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET severity=excluded.severity,state=excluded.state,message=excluded.message,updated_at=excluded.updated_at,acknowledged_at=excluded.acknowledged_at,closed_at=excluded.closed_at`, a.ID, a.TowerID, a.Rule, a.Severity, a.State, a.Message, a.OpenedAt.Format(time.RFC3339Nano), a.UpdatedAt.Format(time.RFC3339Nano), nullableTime(a.AcknowledgedAt), nullableTime(a.ClosedAt))
	return e
}
func nullableTime(v *time.Time) any {
	if v == nil {
		return nil
	}
	return v.Format(time.RFC3339Nano)
}
func (s AlertStore) Get(ctx context.Context, id string) (model.Alert, error) {
	var a model.Alert
	var o, u string
	var ack, cl sql.NullString
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at FROM alerts WHERE id=?`, id).Scan(&a.ID, &a.TowerID, &a.Rule, &a.Severity, &a.State, &a.Message, &o, &u, &ack, &cl)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return a, model.ErrNotFound
		}
		return a, e
	}
	a.OpenedAt, _ = time.Parse(time.RFC3339Nano, o)
	a.UpdatedAt, _ = time.Parse(time.RFC3339Nano, u)
	if ack.Valid && ack.String != "" {
		x, _ := time.Parse(time.RFC3339Nano, ack.String)
		a.AcknowledgedAt = &x
	}
	if cl.Valid && cl.String != "" {
		x, _ := time.Parse(time.RFC3339Nano, cl.String)
		a.ClosedAt = &x
	}
	return a, nil
}
func (s AlertStore) List(ctx context.Context, state string) ([]model.Alert, error) {
	q := `SELECT id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at FROM alerts ORDER BY updated_at DESC`
	args := []any{}
	if state != "" {
		q = `SELECT id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at FROM alerts WHERE state=? ORDER BY updated_at DESC`
		args = append(args, state)
	}
	rows, e := s.DB.SQL.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Alert{}
	for rows.Next() {
		var a model.Alert
		var o, u string
		var ack, cl sql.NullString
		if e := rows.Scan(&a.ID, &a.TowerID, &a.Rule, &a.Severity, &a.State, &a.Message, &o, &u, &ack, &cl); e != nil {
			return nil, e
		}
		a.OpenedAt, _ = time.Parse(time.RFC3339Nano, o)
		a.UpdatedAt, _ = time.Parse(time.RFC3339Nano, u)
		if ack.Valid && ack.String != "" {
			x, _ := time.Parse(time.RFC3339Nano, ack.String)
			a.AcknowledgedAt = &x
		}
		if cl.Valid && cl.String != "" {
			x, _ := time.Parse(time.RFC3339Nano, cl.String)
			a.ClosedAt = &x
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
