package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type RetentionStore struct{ DB *DB }

func (s RetentionStore) Purge(ctx context.Context, p model.RetentionPolicy, now time.Time) (int64, error) {
	var total int64
	if p.ReadingDays > 0 {
		r, e := s.DB.SQL.ExecContext(ctx, `DELETE FROM readings WHERE recorded_at < ?`, now.AddDate(0, 0, -p.ReadingDays).Format(time.RFC3339Nano))
		if e != nil {
			return total, e
		}
		n, _ := r.RowsAffected()
		total += n
	}
	if p.AuditDays > 0 {
		r, e := s.DB.SQL.ExecContext(ctx, `DELETE FROM audit_events WHERE created_at < ?`, now.AddDate(0, 0, -p.AuditDays).Format(time.RFC3339Nano))
		if e != nil {
			return total, e
		}
		n, _ := r.RowsAffected()
		total += n
	}
	return total, nil
}
