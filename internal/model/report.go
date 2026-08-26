package model

import "time"

type ReportFilter struct {
	TowerID string
	From    time.Time
	To      time.Time
	Sensors []string
}
type ReportRow struct {
	Timestamp  time.Time
	Sensor     string
	Value      float64
	Quality    string
	AlertState string
}
type DiagnosticReport struct {
	ID        string
	Tower     Tower
	Filter    ReportFilter
	Rows      []ReportRow
	Findings  []string
	CreatedAt time.Time
}
type ExportResult struct {
	Format    string
	Filename  string
	Bytes     int
	CreatedAt time.Time
}
type RetentionPolicy struct {
	ReadingDays int
	AuditDays   int
	BatchSize   int
}
