package store

import (
	"context"
	"time"
)

type TowerStats struct {
	TowerID     string
	Readings    int64
	Alerts      int64
	LastReading *time.Time
	OpenAlerts  int64
}

func (s DB) Stats(ctx context.Context, tower string) (TowerStats, error) {
	var x TowerStats
	x.TowerID = tower
	if e := s.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE tower_id=?`, tower).Scan(&x.Readings); e != nil {
		return x, e
	}
	if e := s.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE tower_id=?`, tower).Scan(&x.Alerts); e != nil {
		return x, e
	}
	if e := s.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE tower_id=? AND state<> 'closed'`, tower).Scan(&x.OpenAlerts); e != nil {
		return x, e
	}
	var ts string
	if e := s.SQL.QueryRowContext(ctx, `SELECT recorded_at FROM readings WHERE tower_id=? ORDER BY recorded_at DESC LIMIT 1`, tower).Scan(&ts); e == nil {
		t, _ := time.Parse(time.RFC3339Nano, ts)
		x.LastReading = &t
	}
	return x, nil
}
