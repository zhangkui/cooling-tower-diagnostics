package model

import "time"

type TrendPoint struct {
	At       time.Time
	Value    float64
	Baseline float64
	Delta    float64
}
type Trend struct {
	Sensor     string
	Direction  string
	Slope      float64
	Confidence float64
	Points     []TrendPoint
}
type Anomaly struct {
	ID       string
	TowerID  string
	Sensor   string
	At       time.Time
	Value    float64
	Expected float64
	Score    float64
	Reason   string
	Resolved bool
}
type WaterBalance struct {
	TowerID     string    `json:"tower_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Makeup      float64   `json:"makeup_m3"`
	Blowdown    float64   `json:"blowdown_m3"`
	Evaporation float64   `json:"evaporation_m3"`
	Unaccounted float64   `json:"unaccounted_m3"`
}

// FlowCoverage describes the usable portion of a metered-flow series. Volume
// is only integrated across adjacent samples, never extrapolated into a gap.
type FlowCoverage struct {
	Sensor        string  `json:"sensor"`
	Samples       int     `json:"samples"`
	CoveredSecond float64 `json:"covered_seconds"`
	GapCount      int     `json:"gap_count"`
	Volume        float64 `json:"volume_m3"`
}

type BalanceAssessment struct {
	Balance   WaterBalance   `json:"balance"`
	Status    string         `json:"status"`
	Tolerance float64        `json:"tolerance_m3"`
	Coverage  []FlowCoverage `json:"coverage"`
	Findings  []string       `json:"findings"`
}
type EfficiencyTarget struct {
	TowerID   string
	Sensor    string
	Target    float64
	Tolerance float64
	Active    bool
}
type OperatorNote struct {
	ID        string
	TowerID   string
	AlertID   string
	Author    string
	Body      string
	CreatedAt time.Time
}

type Capacity struct {
	TowerID     string
	RatedFlow   float64
	CurrentFlow float64
	Utilization float64
	Headroom    float64
	Status      string
}

type Forecast struct {
	Sensor         string
	HorizonMinutes int
	Predicted      []float64
	Lower          []float64
	Upper          []float64
	GeneratedAt    time.Time
}
