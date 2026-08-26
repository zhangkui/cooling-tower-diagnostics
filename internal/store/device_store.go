package store

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"database/sql"
	"errors"
	"time"
)

type DeviceStore struct{ DB *DB }

func (s DeviceStore) Put(ctx context.Context, d model.Device) error {
	_, err := s.DB.SQL.ExecContext(ctx, `INSERT INTO devices(id,tower_id,protocol,address,state,last_heartbeat,failure) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET protocol=excluded.protocol,address=excluded.address,state=excluded.state,last_heartbeat=excluded.last_heartbeat,failure=excluded.failure`, d.ID, d.TowerID, d.Protocol, d.Address, d.State, d.LastHeartbeat.Format(time.RFC3339Nano), d.Failure)
	return err
}

func (s DeviceStore) Get(ctx context.Context, id string) (model.Device, error) {
	var d model.Device
	var heartbeat string
	err := s.DB.SQL.QueryRowContext(ctx, `SELECT id,tower_id,protocol,address,state,last_heartbeat,failure FROM devices WHERE id=?`, id).Scan(&d.ID, &d.TowerID, &d.Protocol, &d.Address, &d.State, &heartbeat, &d.Failure)
	if errors.Is(err, sql.ErrNoRows) {
		return d, model.ErrNotFound
	}
	if err != nil {
		return d, err
	}
	d.LastHeartbeat, err = time.Parse(time.RFC3339Nano, heartbeat)
	return d, err
}

func (s DeviceStore) ForTower(ctx context.Context, towerID string) ([]model.Device, error) {
	rows, err := s.DB.SQL.QueryContext(ctx, `SELECT id,tower_id,protocol,address,state,last_heartbeat,failure FROM devices WHERE tower_id=? ORDER BY last_heartbeat DESC`, towerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	devices := []model.Device{}
	for rows.Next() {
		var d model.Device
		var heartbeat string
		if err := rows.Scan(&d.ID, &d.TowerID, &d.Protocol, &d.Address, &d.State, &heartbeat, &d.Failure); err != nil {
			return nil, err
		}
		if d.LastHeartbeat, err = time.Parse(time.RFC3339Nano, heartbeat); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (s DeviceStore) UpdateState(ctx context.Context, id, state, failure string, at time.Time) error {
	result, err := s.DB.SQL.ExecContext(ctx, `UPDATE devices SET state=?,failure=?,last_heartbeat=? WHERE id=?`, state, failure, at.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}
