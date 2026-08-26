package model

import "time"

type CreateTowerRequest struct {
	Name     string  `json:"name"`
	Site     string  `json:"site"`
	DesignKW float64 `json:"design_kw"`
}
type CreateSensorRequest struct {
	TowerID string  `json:"tower_id"`
	Kind    string  `json:"kind"`
	Unit    string  `json:"unit"`
	Offset  float64 `json:"offset"`
	Scale   float64 `json:"scale"`
}
type CreateDeviceRequest struct {
	TowerID  string `json:"tower_id"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}
type HeartbeatRequest struct {
	Failure string `json:"failure"`
}
type ReadingRequest struct {
	TowerID    string     `json:"tower_id"`
	Sensor     string     `json:"sensor"`
	Value      float64    `json:"value"`
	Unit       string     `json:"unit"`
	RecordedAt *time.Time `json:"recorded_at,omitempty"`
}
type ReplayRequest struct {
	TowerID     string           `json:"tower_id"`
	Frames      []ReadingRequest `json:"frames"`
	DelayMillis int              `json:"delay_millis"`
}
type ThresholdRequest struct {
	Sensor        string  `json:"sensor"`
	Warning       float64 `json:"warning"`
	Critical      float64 `json:"critical"`
	Direction     string  `json:"direction"`
	WindowMinutes int     `json:"window_minutes"`
}
type StateRequest struct {
	Actor string `json:"actor"`
	Note  string `json:"note"`
}
type DiagnosisRequest struct {
	TowerID      string `json:"tower_id"`
	SinceMinutes int    `json:"since_minutes"`
}

type WaterBalanceRequest struct {
	TowerID       string    `json:"tower_id"`
	From          time.Time `json:"from"`
	To            time.Time `json:"to"`
	Tolerance     float64   `json:"tolerance_m3"`
	MaxGapMinutes int       `json:"max_gap_minutes"`
}

type CreateMaintenanceRequest struct {
	TowerID         string     `json:"tower_id"`
	Kind            string     `json:"kind"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	DurationMinutes int        `json:"duration_minutes"`
	Owner           string     `json:"owner"`
}
