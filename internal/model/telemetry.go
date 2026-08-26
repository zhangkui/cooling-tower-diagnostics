package model

import "time"

type TelemetryFrame struct {
	DeviceID string
	Sequence int64
	SentAt   time.Time
	Fields   map[string]float64
	Tags     map[string]string
}
type ProtocolStatus struct {
	DeviceID  string
	Connected bool
	LastSeen  time.Time
	Frames    int64
	Errors    int64
}
type CalibrationProfile struct {
	Sensor      string
	Version     string
	EffectiveAt time.Time
	Offset      float64
	Scale       float64
	Min         float64
	Max         float64
}
type QualityEvent struct {
	ReadingID string
	Code      string
	Message   string
	Severity  string
	CreatedAt time.Time
}
type MaintenancePlan struct {
	ID              string
	TowerID         string
	Kind            string
	DueAt           time.Time
	DurationMinutes int
	Owner           string
	Status          string
}
