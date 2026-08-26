package model

import "time"

type DataGap struct {
	TowerID  string
	Sensor   string
	From     time.Time
	To       time.Time
	Expected int
	Actual   int
	Missing  int
}
type ThresholdBreach struct {
	TowerID  string
	Sensor   string
	Severity string
	Value    float64
	Limit    float64
	At       time.Time
}
type EnergySample struct {
	TowerID    string
	At         time.Time
	InputKW    float64
	OutputKW   float64
	Efficiency float64
}
