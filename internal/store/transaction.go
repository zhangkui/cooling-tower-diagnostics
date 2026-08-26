package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"fmt"
)

type Tx struct {
	DB  *DB
	SQL *sql.Tx
}

func (s *DB) Begin(ctx context.Context) (Tx, error) {
	tx, e := s.SQL.BeginTx(ctx, nil)
	return Tx{DB: s, SQL: tx}, e
}
func (t Tx) Commit() error   { return t.SQL.Commit() }
func (t Tx) Rollback() error { return t.SQL.Rollback() }
func (t Tx) InsertAudit(ctx context.Context, e model.AuditEvent) error {
	_, err := t.SQL.ExecContext(ctx, `INSERT INTO audit_events(id,entity,entity_id,action,details,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.Entity, e.EntityID, e.Action, e.Details, e.CreatedAt)
	if err != nil {
		return fmt.Errorf("transaction audit: %w", err)
	}
	return nil
}
func (t Tx) InsertReading(ctx context.Context, r model.Reading) error {
	_, err := t.SQL.ExecContext(ctx, `INSERT INTO readings(id,tower_id,sensor,value,raw_value,unit,quality,recorded_at) VALUES(?,?,?,?,?,?,?,?)`, r.ID, r.TowerID, r.Sensor, r.Value, r.RawValue, r.Unit, r.Quality, r.RecordedAt)
	return err
}
