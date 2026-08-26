package service

import (
	"context"
	"cooling-tower-diagnostics/internal/alert"
	"cooling-tower-diagnostics/internal/clock"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/ingest"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/rules"
	"cooling-tower-diagnostics/internal/store"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var readingSequence atomic.Uint64
var auditSequence atomic.Uint64

type Service struct {
	Towers      store.TowerStore
	Readings    store.ReadingStore
	Alerts      store.AlertStore
	Audits      store.AuditStore
	Thresholds  store.ThresholdStore
	Engine      diagnostics.Engine
	Clock       clock.Clock
	Decoder     ingest.Decoder
	Sensors     store.SensorStore
	Maintenance store.MaintenanceStore
	Queries     store.QueryStore
	Devices     store.DeviceStore
}

func (s *Service) CreateTower(ctx context.Context, r model.CreateTowerRequest) (model.Tower, error) {
	if strings.TrimSpace(r.Name) == "" || r.DesignKW <= 0 {
		return model.Tower{}, model.ErrInvalid
	}
	now := s.Clock.Now()
	t := model.Tower{ID: fmt.Sprintf("tower-%d", now.UnixNano()), Name: strings.TrimSpace(r.Name), Site: strings.TrimSpace(r.Site), DesignKW: r.DesignKW, Status: model.TowerOnline, CreatedAt: now}
	if err := s.Towers.Create(ctx, t); err != nil {
		return t, err
	}
	_ = s.audit(ctx, "tower", t.ID, "created", t.Name)
	return t, nil
}
func (s *Service) ListTowers(ctx context.Context) ([]model.Tower, error) { return s.Towers.List(ctx) }
func (s *Service) SetTowerStatus(ctx context.Context, id, status string) error {
	if status != model.TowerOnline && status != model.TowerPaused {
		return model.ErrInvalid
	}
	if err := s.Towers.SetStatus(ctx, id, status); err != nil {
		return err
	}
	return s.audit(ctx, "tower", id, "status", status)
}
func (s *Service) ConfigureThreshold(ctx context.Context, r model.ThresholdRequest) error {
	if r.Sensor == "" || r.WindowMinutes <= 0 {
		return model.ErrInvalid
	}
	return s.Thresholds.Put(ctx, model.Threshold{Sensor: r.Sensor, Warning: r.Warning, Critical: r.Critical, Direction: r.Direction, WindowMinutes: r.WindowMinutes})
}
func (s *Service) ListThresholds(ctx context.Context) ([]model.Threshold, error) {
	return s.Thresholds.List(ctx)
}
func (s *Service) IngestReading(ctx context.Context, r model.ReadingRequest) (model.Reading, error) {
	if r.TowerID == "" || r.Sensor == "" {
		return model.Reading{}, model.ErrInvalid
	}
	now := s.Clock.Now()
	if r.RecordedAt != nil {
		now = *r.RecordedAt
	}
	v := r.Value
	if sensors, err := s.Sensors.ForTower(ctx, r.TowerID); err == nil {
		for _, sensor := range sensors {
			if sensor.Kind == strings.ToLower(strings.TrimSpace(r.Sensor)) || sensor.Enabled {
				v = v*sensor.Scale + sensor.Offset
				break
			}
		}
	}
	quality := model.ReadingGood
	sensor := strings.ToLower(strings.TrimSpace(r.Sensor))
	reading := model.Reading{ID: fmt.Sprintf("reading-%d-%d", now.UnixNano(), readingSequence.Add(1)), TowerID: r.TowerID, Sensor: sensor, RawValue: r.Value, Value: v, Unit: r.Unit, Quality: quality, RecordedAt: now}
	if err := s.Readings.Add(ctx, reading); err != nil {
		return reading, err
	}
	return reading, nil
}
func (s *Service) Replay(ctx context.Context, r model.ReplayRequest) error {
	return ingest.Replay(ctx, r.Frames, time.Duration(r.DelayMillis)*time.Millisecond, func(c context.Context, x model.ReadingRequest) error { _, e := s.IngestReading(c, x); return e })
}
func (s *Service) Diagnose(ctx context.Context, r model.DiagnosisRequest) (model.Diagnosis, error) {
	d, e := s.Engine.Run(ctx, r.TowerID, r.SinceMinutes)
	if e != nil {
		return d, e
	}
	d.Findings = diagnostics.Explain(d)
	return d, nil
}
func (s *Service) EvaluateAndAlert(ctx context.Context, r model.DiagnosisRequest) (model.Diagnosis, error) {
	d, e := s.Diagnose(ctx, r)
	if e != nil {
		return d, e
	}
	for _, f := range d.Findings {
		if strings.HasPrefix(f, "critical:") {
			a := alert.Open(s.Clock.Now(), r.TowerID, "quality-window", model.SeverityCritical, f)
			if err := s.Alerts.Upsert(ctx, a); err != nil {
				return d, err
			}
			_ = s.audit(ctx, "alert", a.ID, "opened", f)
		}
	}
	return d, nil
}
func (s *Service) ListAlerts(ctx context.Context, state string) ([]model.Alert, error) {
	return s.Alerts.List(ctx, state)
}
func (s *Service) TransitionAlert(ctx context.Context, id, action string) error {
	a, e := s.Alerts.Get(ctx, id)
	if e != nil {
		return e
	}
	now := s.Clock.Now()
	switch action {
	case "ack":
		e = alert.Acknowledge(&a, now)
	case "escalate":
		e = alert.Escalate(&a, now)
	case "close":
		e = alert.Close(&a, now)
	default:
		return model.ErrInvalid
	}
	if e != nil {
		return e
	}
	if e = s.Alerts.Upsert(ctx, a); e != nil {
		return e
	}
	return s.audit(ctx, "alert", id, action, "")
}
func (s *Service) ListAudit(ctx context.Context, entity string) ([]model.AuditEvent, error) {
	return s.Audits.List(ctx, entity)
}
func (s *Service) audit(ctx context.Context, entity, id, action, details string) error {
	b, _ := json.Marshal(map[string]string{"details": details})
	now := s.Clock.Now()
	return s.Audits.Append(ctx, model.AuditEvent{ID: fmt.Sprintf("audit-%d-%d", now.UnixNano(), auditSequence.Add(1)), Entity: entity, EntityID: id, Action: action, Details: string(b), CreatedAt: now})
}
func (s *Service) Seed(ctx context.Context) error {
	ts, _ := s.Towers.List(ctx)
	if len(ts) > 0 {
		return nil
	}
	t, _ := s.CreateTower(ctx, model.CreateTowerRequest{Name: "北区一号冷却塔", Site: "North", DesignKW: 420})
	_ = s.Thresholds.Put(ctx, model.Threshold{Sensor: "conductivity", Warning: 1500, Critical: 2200, Direction: "above", WindowMinutes: 30})
	_ = s.Thresholds.Put(ctx, model.Threshold{Sensor: "outlet_temp", Warning: 34, Critical: 38, Direction: "above", WindowMinutes: 15})
	_ = s.Sensors.Put(ctx, model.Sensor{ID: "sensor-" + t.ID + "-conductivity", TowerID: t.ID, Kind: "conductivity", Unit: "uS/cm", Scale: 1, Enabled: true})
	_ = s.Sensors.Put(ctx, model.Sensor{ID: "sensor-" + t.ID + "-outlet-temp", TowerID: t.ID, Kind: "outlet_temp", Unit: "C", Scale: 1, Enabled: true})
	_, _ = s.IngestReading(ctx, model.ReadingRequest{TowerID: t.ID, Sensor: "conductivity", Value: 1200, Unit: "uS/cm"})
	return nil
}

func (s *Service) ListSensors(ctx context.Context, towerID string) ([]model.Sensor, error) {
	return s.Sensors.ForTower(ctx, towerID)
}

func (s *Service) SaveSensor(ctx context.Context, req model.CreateSensorRequest) (model.Sensor, error) {
	if req.TowerID == "" || req.Kind == "" || req.Unit == "" {
		return model.Sensor{}, model.ErrInvalid
	}
	if req.Scale == 0 {
		req.Scale = 1
	}
	sensor := model.Sensor{ID: fmt.Sprintf("sensor-%d", s.Clock.Now().UnixNano()), TowerID: req.TowerID, Kind: req.Kind, Unit: req.Unit, Offset: req.Offset, Scale: req.Scale, Enabled: true}
	if err := s.Sensors.Put(ctx, sensor); err != nil {
		return sensor, err
	}
	_ = s.audit(ctx, "sensor", sensor.ID, "saved", sensor.Kind)
	return sensor, nil
}

func (s *Service) PlanMaintenance(ctx context.Context, req model.CreateMaintenanceRequest) (model.MaintenancePlan, error) {
	if req.TowerID == "" || req.Kind == "" {
		return model.MaintenancePlan{}, model.ErrInvalid
	}
	due := s.Clock.Now().Add(24 * time.Hour)
	if req.DueAt != nil {
		due = req.DueAt.UTC()
	}
	plan := model.MaintenancePlan{ID: fmt.Sprintf("maintenance-%d", s.Clock.Now().UnixNano()), TowerID: req.TowerID, Kind: req.Kind, DueAt: due, DurationMinutes: req.DurationMinutes, Owner: req.Owner, Status: "planned"}
	if plan.DurationMinutes <= 0 {
		plan.DurationMinutes = 60
	}
	if err := s.Maintenance.Put(ctx, plan); err != nil {
		return plan, err
	}
	_ = s.audit(ctx, "maintenance", plan.ID, "planned", plan.Kind)
	return plan, nil
}

func (s *Service) ListMaintenance(ctx context.Context, towerID string) ([]model.MaintenancePlan, error) {
	return s.Maintenance.List(ctx, towerID)
}

func (s *Service) CompleteMaintenance(ctx context.Context, id string) error {
	if err := s.Maintenance.Complete(ctx, id); err != nil {
		return err
	}
	return s.audit(ctx, "maintenance", id, "completed", "")
}

func (s *Service) TowerStats(ctx context.Context, towerID string) (store.TowerStats, error) {
	return s.Queries.DB.Stats(ctx, towerID)
}

func (s *Service) Report(ctx context.Context, filter model.ReportFilter) (model.DiagnosticReport, error) {
	reporter := diagnostics.Reporter{Readings: s.Readings, Engine: s.Engine}
	return reporter.Build(ctx, filter)
}

func (s *Service) Forecast(ctx context.Context, towerID, sensor string, horizon int) (model.Forecast, error) {
	forecaster := diagnostics.Forecaster{Readings: s.Readings, Model: rules.ForecastModel{Window: 8, Decay: 0.8, TrendWeight: 1}}
	return forecaster.Run(ctx, towerID, sensor, horizon)
}

func (s *Service) Capacity(ctx context.Context, towerID string, ratedFlow float64) (model.Capacity, error) {
	checker := diagnostics.CapacityService{Readings: s.Readings, Towers: s.Towers}
	return checker.Evaluate(ctx, towerID, ratedFlow)
}
func (s *Service) ParsePayload(b []byte) (model.ReadingRequest, error) { return s.Decoder.Decode(b) }
func (s *Service) Summary(ctx context.Context, tower string) (map[string]any, error) {
	n, e := s.Readings.Count(ctx, tower)
	if e != nil {
		return nil, e
	}
	last, e := s.Readings.Last(ctx, tower)
	if e != nil {
		return map[string]any{"count": n}, nil
	}
	return map[string]any{"count": n, "last": last}, nil
}
