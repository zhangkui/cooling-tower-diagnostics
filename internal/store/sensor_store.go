package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
)

type SensorStore struct{ DB *DB }

func (s SensorStore) Put(ctx context.Context, x model.Sensor) error {
	_, e := s.DB.SQL.ExecContext(ctx, `INSERT INTO sensors(id,tower_id,kind,unit,offset,scale,enabled) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET unit=excluded.unit,offset=excluded.offset,scale=excluded.scale,enabled=excluded.enabled`, x.ID, x.TowerID, x.Kind, x.Unit, x.Offset, x.Scale, x.Enabled)
	return e
}
func (s SensorStore) Get(ctx context.Context, id string) (model.Sensor, error) {
	var x model.Sensor
	var enabled int
	e := s.DB.SQL.QueryRowContext(ctx, `SELECT id,tower_id,kind,unit,offset,scale,enabled FROM sensors WHERE id=?`, id).Scan(&x.ID, &x.TowerID, &x.Kind, &x.Unit, &x.Offset, &x.Scale, &enabled)
	if e == sql.ErrNoRows {
		return x, model.ErrNotFound
	}
	x.Enabled = enabled != 0
	return x, e
}
func (s SensorStore) ForTower(ctx context.Context, tower string) ([]model.Sensor, error) {
	rows, e := s.DB.SQL.QueryContext(ctx, `SELECT id,tower_id,kind,unit,offset,scale,enabled FROM sensors WHERE tower_id=? ORDER BY kind`, tower)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Sensor{}
	for rows.Next() {
		var x model.Sensor
		var en int
		if e := rows.Scan(&x.ID, &x.TowerID, &x.Kind, &x.Unit, &x.Offset, &x.Scale, &en); e != nil {
			return nil, e
		}
		x.Enabled = en != 0
		out = append(out, x)
	}
	return out, rows.Err()
}
