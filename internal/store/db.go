package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type DB struct{ SQL *sql.DB }

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepathDir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &DB{SQL: db}, nil
}

func filepathDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			if i == 0 {
				return path[:1]
			}
			return path[:i]
		}
	}
	return "."
}

func (d *DB) Close() error { return d.SQL.Close() }

func (d *DB) Migrate(ctx context.Context) error {
	stmts := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS towers (id TEXT PRIMARY KEY, name TEXT NOT NULL, site TEXT NOT NULL, design_kw REAL NOT NULL, status TEXT NOT NULL, created_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS sensors (id TEXT PRIMARY KEY, tower_id TEXT NOT NULL REFERENCES towers(id), kind TEXT NOT NULL, unit TEXT NOT NULL, offset REAL NOT NULL, scale REAL NOT NULL, enabled INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS readings (id TEXT PRIMARY KEY, tower_id TEXT NOT NULL REFERENCES towers(id), sensor TEXT NOT NULL, value REAL NOT NULL, raw_value REAL NOT NULL, unit TEXT NOT NULL, quality TEXT NOT NULL, recorded_at TEXT NOT NULL)`,
		`CREATE INDEX IF NOT EXISTS idx_readings_tower_sensor_time ON readings(tower_id, sensor, recorded_at)`,
		`CREATE TABLE IF NOT EXISTS thresholds (sensor TEXT PRIMARY KEY, warning REAL NOT NULL, critical REAL NOT NULL, direction TEXT NOT NULL, window_minutes INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS alerts (id TEXT PRIMARY KEY, tower_id TEXT NOT NULL, rule TEXT NOT NULL, severity TEXT NOT NULL, state TEXT NOT NULL, message TEXT NOT NULL, opened_at TEXT NOT NULL, updated_at TEXT NOT NULL, acknowledged_at TEXT, closed_at TEXT)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_alert_open ON alerts(tower_id, rule, state) WHERE state <> 'closed'`,
		`CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY, entity TEXT NOT NULL, entity_id TEXT NOT NULL, action TEXT NOT NULL, details TEXT NOT NULL, created_at TEXT NOT NULL)`,
	}
	for _, stmt := range stmts {
		if _, err := d.SQL.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
