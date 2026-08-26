package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"time"
)

type NoteStore struct{ DB *DB }

func (s NoteStore) Add(ctx context.Context, n model.OperatorNote) error {
	_, e := s.DB.SQL.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS operator_notes(id TEXT PRIMARY KEY,tower_id TEXT,alert_id TEXT,author TEXT,body TEXT,created_at TEXT)`)
	if e != nil {
		return e
	}
	_, e = s.DB.SQL.ExecContext(ctx, `INSERT INTO operator_notes(id,tower_id,alert_id,author,body,created_at) VALUES(?,?,?,?,?,?)`, n.ID, n.TowerID, n.AlertID, n.Author, n.Body, n.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (s NoteStore) List(ctx context.Context, tower string) ([]model.OperatorNote, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT id,tower_id,alert_id,author,body,created_at FROM operator_notes WHERE tower_id=? ORDER BY created_at DESC`, tower)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.OperatorNote{}
	for rows.Next() {
		var n model.OperatorNote
		var ts string
		if e := rows.Scan(&n.ID, &n.TowerID, &n.AlertID, &n.Author, &n.Body, &ts); e != nil {
			return nil, e
		}
		n.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, n)
	}
	return out, rows.Err()
}
