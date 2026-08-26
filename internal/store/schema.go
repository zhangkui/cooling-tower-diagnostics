package store

import "context"

func (d *DB) EnsureOptionalSchema(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS calibrations(sensor TEXT PRIMARY KEY,version TEXT NOT NULL,effective_at TEXT NOT NULL,offset REAL NOT NULL,scale REAL NOT NULL,min_value REAL NOT NULL,max_value REAL NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS maintenance_plans(id TEXT PRIMARY KEY,tower_id TEXT NOT NULL,kind TEXT NOT NULL,due_at TEXT NOT NULL,duration_minutes INTEGER NOT NULL,owner TEXT NOT NULL,status TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS anomalies(id TEXT PRIMARY KEY,tower_id TEXT NOT NULL,sensor TEXT NOT NULL,at TEXT NOT NULL,value REAL NOT NULL,expected REAL NOT NULL,score REAL NOT NULL,reason TEXT NOT NULL,resolved INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS operator_notes(id TEXT PRIMARY KEY,tower_id TEXT NOT NULL,alert_id TEXT,author TEXT NOT NULL,body TEXT NOT NULL,created_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS devices(id TEXT PRIMARY KEY,tower_id TEXT NOT NULL REFERENCES towers(id),protocol TEXT NOT NULL,address TEXT NOT NULL,state TEXT NOT NULL,last_heartbeat TEXT NOT NULL,failure TEXT NOT NULL)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_tower ON devices(tower_id)`,
	}
	for _, statement := range statements {
		if _, err := d.SQL.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
