package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TowerStore struct{ DB *DB }

func (s TowerStore) Create(ctx context.Context, t model.Tower) error {
	_, err := s.DB.SQL.ExecContext(ctx, `INSERT INTO towers(id,name,site,design_kw,status,created_at) VALUES(?,?,?,?,?,?)`, t.ID, t.Name, t.Site, t.DesignKW, t.Status, t.CreatedAt.Format(time.RFC3339Nano))
	return err
}
func (s TowerStore) Get(ctx context.Context, id string) (model.Tower, error) {
	var t model.Tower
	var ts string
	err := s.DB.SQL.QueryRowContext(ctx, `SELECT id,name,site,design_kw,status,created_at FROM towers WHERE id=?`, id).Scan(&t.ID, &t.Name, &t.Site, &t.DesignKW, &t.Status, &ts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, model.ErrNotFound
		}
		return t, fmt.Errorf("get tower: %w", err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return t, nil
}
func (s TowerStore) List(ctx context.Context) ([]model.Tower, error) {
	rows, err := s.DB.SQL.QueryContext(ctx, `SELECT id,name,site,design_kw,status,created_at FROM towers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Tower, 0)
	for rows.Next() {
		var t model.Tower
		var ts string
		if err := rows.Scan(&t.ID, &t.Name, &t.Site, &t.DesignKW, &t.Status, &ts); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, t)
	}
	return out, rows.Err()
}
func (s TowerStore) SetStatus(ctx context.Context, id, status string) error {
	result, err := s.DB.SQL.ExecContext(ctx, `UPDATE towers SET status=? WHERE id=?`, status, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}
