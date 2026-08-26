package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"fmt"
	"strings"
	"time"
)

var balanceSensors = []string{"makeup_flow", "blowdown_flow", "evaporation_flow"}

// AssessWaterBalance turns durable meter samples into a time-bounded water
// account. Readings are flow rates in m3/h; the rule layer integrates them.
func (s *Service) AssessWaterBalance(ctx context.Context, request model.WaterBalanceRequest) (model.BalanceAssessment, error) {
	request.TowerID = strings.TrimSpace(request.TowerID)
	if request.TowerID == "" || request.From.IsZero() || request.To.IsZero() || !request.To.After(request.From) || request.Tolerance < 0 {
		return model.BalanceAssessment{}, model.ErrInvalid
	}
	if _, err := s.Towers.Get(ctx, request.TowerID); err != nil {
		return model.BalanceAssessment{}, err
	}
	if request.MaxGapMinutes <= 0 {
		request.MaxGapMinutes = 15
	}
	if request.Tolerance == 0 {
		request.Tolerance = 0.5
	}

	volumes := make(map[string]float64, len(balanceSensors))
	coverage := make([]model.FlowCoverage, 0, len(balanceSensors))
	for _, sensor := range balanceSensors {
		rows, err := s.Queries.ReadingsByRange(ctx, model.ReportFilter{TowerID: request.TowerID, From: request.From, To: request.To, Sensors: []string{sensor}})
		if err != nil {
			return model.BalanceAssessment{}, err
		}
		item := rules.IntegrateFlow(sensor, rows, time.Duration(request.MaxGapMinutes)*time.Minute)
		coverage = append(coverage, item)
		volumes[sensor] = item.Volume
	}

	balance := rules.WaterBalance(request.TowerID, volumes["makeup_flow"], volumes["blowdown_flow"], volumes["evaporation_flow"])
	balance.PeriodStart, balance.PeriodEnd = request.From.UTC(), request.To.UTC()
	assessment := model.BalanceAssessment{Balance: balance, Status: rules.BalanceStatus(balance, request.Tolerance), Tolerance: request.Tolerance, Coverage: coverage}
	for _, item := range coverage {
		if item.Samples < 2 {
			assessment.Findings = append(assessment.Findings, fmt.Sprintf("%s has insufficient samples", item.Sensor))
		}
		if item.GapCount > 0 {
			assessment.Findings = append(assessment.Findings, fmt.Sprintf("%s excluded %d collection gaps", item.Sensor, item.GapCount))
		}
	}
	if assessment.Status != "balanced" {
		assessment.Findings = append(assessment.Findings, "water account requires operator review")
	}
	if err := s.audit(ctx, "water-balance", request.TowerID, "assessed", assessment.Status); err != nil {
		return model.BalanceAssessment{}, err
	}
	return assessment, nil
}
