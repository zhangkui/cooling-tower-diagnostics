package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type AuditStore struct{ DB *DB }

func (s AuditStore) Append(ctx context.Context, e model.AuditEvent) error {
	_, x := s.DB.SQL.ExecContext(ctx, `INSERT INTO audit_events(id,entity,entity_id,action,details,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.Entity, e.EntityID, e.Action, e.Details, e.CreatedAt.Format(time.RFC3339Nano))
	return x
}
func (s AuditStore) List(ctx context.Context, entity string) ([]model.AuditEvent, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT id,entity,entity_id,action,details,created_at FROM audit_events WHERE entity=? ORDER BY created_at DESC`, entity)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.AuditEvent{}
	for rows.Next() {
		var a model.AuditEvent
		var ts string
		if e := rows.Scan(&a.ID, &a.Entity, &a.EntityID, &a.Action, &a.Details, &ts); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, a)
	}
	return out, rows.Err()
}
