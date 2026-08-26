package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type MaintenanceStore struct{ DB *DB }

func (s MaintenanceStore) Put(ctx context.Context, p model.MaintenancePlan) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO maintenance_plans(id,tower_id,kind,due_at,duration_minutes,owner,status) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET due_at=excluded.due_at,owner=excluded.owner,status=excluded.status`, p.ID, p.TowerID, p.Kind, p.DueAt.Format(time.RFC3339Nano), p.DurationMinutes, p.Owner, p.Status)
	return e
}
func (s MaintenanceStore) List(ctx context.Context, tower string) ([]model.MaintenancePlan, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT id,tower_id,kind,due_at,duration_minutes,owner,status FROM maintenance_plans WHERE tower_id=? ORDER BY due_at`, tower)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.MaintenancePlan{}
	for rows.Next() {
		var p model.MaintenancePlan
		var ts string
		if e := rows.Scan(&p.ID, &p.TowerID, &p.Kind, &ts, &p.DurationMinutes, &p.Owner, &p.Status); e != nil {
			return nil, e
		}
		p.DueAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, p)
	}
	return out, rows.Err()
}
func (s MaintenanceStore) Complete(ctx context.Context, id string) error {
	_, e := s.DB.SQL.ExecContext(ctx, `UPDATE maintenance_plans SET status='done' WHERE id=?`, id)
	return e
}
