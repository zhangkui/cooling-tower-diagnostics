package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"fmt"
	"time"
)

type PaginationStore struct{ DB *DB }

func (s PaginationStore) Towers(ctx context.Context, page, size int) (model.Page, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	var total int
	if e := s.DB.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM towers`).Scan(&total); e != nil {
		return model.Page{}, e
	}
	rows, e := s.DB.SQL.QueryContext(ctx, fmt.Sprintf(`SELECT id,name,site,design_kw,status,created_at FROM towers ORDER BY name LIMIT %d OFFSET %d`, size, (page-1)*size))
	if e != nil {
		return model.Page{}, e
	}
	defer rows.Close()
	items := []model.Tower{}
	for rows.Next() {
		var t model.Tower
		var ts string
		if e := rows.Scan(&t.ID, &t.Name, &t.Site, &t.DesignKW, &t.Status, &ts); e != nil {
			return model.Page{}, e
		}
		items = append(items, t)
	}
	return model.NewPage(page, size, total, items), rows.Err()
}
func (s PaginationStore) Alerts(ctx context.Context, page, size int) (model.Page, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	var total int
	if e := s.DB.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts`).Scan(&total); e != nil {
		return model.Page{}, e
	}
	rows, e := s.DB.SQL.QueryContext(ctx, fmt.Sprintf(`SELECT id,tower_id,rule,severity,state,message,opened_at,updated_at,acknowledged_at,closed_at FROM alerts ORDER BY updated_at DESC LIMIT %d OFFSET %d`, size, (page-1)*size))
	if e != nil {
		return model.Page{}, e
	}
	defer rows.Close()
	items := []model.Alert{}
	for rows.Next() {
		var a model.Alert
		var opened, updated string
		var ack, closed sql.NullString
		if e := rows.Scan(&a.ID, &a.TowerID, &a.Rule, &a.Severity, &a.State, &a.Message, &opened, &updated, &ack, &closed); e != nil {
			return model.Page{}, e
		}
		a.OpenedAt, _ = time.Parse(time.RFC3339Nano, opened)
		a.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
		if ack.Valid && ack.String != "" {
			x, _ := time.Parse(time.RFC3339Nano, ack.String)
			a.AcknowledgedAt = &x
		}
		if closed.Valid && closed.String != "" {
			x, _ := time.Parse(time.RFC3339Nano, closed.String)
			a.ClosedAt = &x
		}
		items = append(items, a)
	}
	return model.NewPage(page, size, total, items), rows.Err()
}
