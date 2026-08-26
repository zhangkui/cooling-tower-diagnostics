package model

import "time"

type Tower struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Site      string    `json:"site"`
	DesignKW  float64   `json:"design_kw"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Sensor struct {
	ID      string  `json:"id"`
	TowerID string  `json:"tower_id"`
	Kind    string  `json:"kind"`
	Unit    string  `json:"unit"`
	Offset  float64 `json:"offset"`
	Scale   float64 `json:"scale"`
	Enabled bool    `json:"enabled"`
}

// Device is the field gateway that submits telemetry for one cooling tower.
// Its heartbeat is persisted so an outage remains visible after a restart.
type Device struct {
	ID            string    `json:"id"`
	TowerID       string    `json:"tower_id"`
	Protocol      string    `json:"protocol"`
	Address       string    `json:"address"`
	State         string    `json:"state"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Failure       string    `json:"failure,omitempty"`
}

type Reading struct {
	ID         string    `json:"id"`
	TowerID    string    `json:"tower_id"`
	Sensor     string    `json:"sensor"`
	Value      float64   `json:"value"`
	RawValue   float64   `json:"raw_value"`
	Unit       string    `json:"unit"`
	Quality    string    `json:"quality"`
	RecordedAt time.Time `json:"recorded_at"`
}

type WindowSummary struct {
	TowerID string    `json:"tower_id"`
	Sensor  string    `json:"sensor"`
	Count   int       `json:"count"`
	Average float64   `json:"average"`
	Minimum float64   `json:"minimum"`
	Maximum float64   `json:"maximum"`
	Spread  float64   `json:"spread"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

type Alert struct {
	ID             string     `json:"id"`
	TowerID        string     `json:"tower_id"`
	Rule           string     `json:"rule"`
	Severity       string     `json:"severity"`
	State          string     `json:"state"`
	Message        string     `json:"message"`
	OpenedAt       time.Time  `json:"opened_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
}

type AuditEvent struct {
	ID        string    `json:"id"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type Diagnosis struct {
	TowerID     string        `json:"tower_id"`
	Window      WindowSummary `json:"window"`
	Findings    []string      `json:"findings"`
	Risk        string        `json:"risk"`
	EnergyIndex float64       `json:"energy_index"`
	GeneratedAt time.Time     `json:"generated_at"`
}

type Threshold struct {
	Sensor        string  `json:"sensor"`
	Warning       float64 `json:"warning"`
	Critical      float64 `json:"critical"`
	Direction     string  `json:"direction"`
	WindowMinutes int     `json:"window_minutes"`
}
