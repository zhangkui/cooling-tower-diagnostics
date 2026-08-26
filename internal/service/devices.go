package service

import (
	"context"
	"cooling-tower-diagnostics/internal/alert"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"strings"
	"time"
)

const deviceHeartbeatMaxAge = 5 * time.Minute

func (s *Service) RegisterDevice(ctx context.Context, request model.CreateDeviceRequest) (model.Device, error) {
	request.TowerID = strings.TrimSpace(request.TowerID)
	request.Protocol = strings.ToLower(strings.TrimSpace(request.Protocol))
	request.Address = strings.TrimSpace(request.Address)
	if request.TowerID == "" || request.Protocol == "" || request.Address == "" {
		return model.Device{}, model.ErrInvalid
	}
	if _, err := s.Towers.Get(ctx, request.TowerID); err != nil {
		return model.Device{}, err
	}
	now := s.Clock.Now()
	device := model.Device{ID: fmt.Sprintf("device-%d", now.UnixNano()), TowerID: request.TowerID, Protocol: request.Protocol, Address: request.Address, State: "connected", LastHeartbeat: now}
	if err := s.Devices.Put(ctx, device); err != nil {
		return model.Device{}, err
	}
	if err := s.audit(ctx, "device", device.ID, "registered", device.Protocol+":"+device.Address); err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (s *Service) HeartbeatDevice(ctx context.Context, id, failure string) (model.Device, error) {
	device, err := s.Devices.Get(ctx, id)
	if err != nil {
		return model.Device{}, err
	}
	state := "connected"
	if strings.TrimSpace(failure) != "" {
		state = "degraded"
	}
	now := s.Clock.Now()
	if err := s.Devices.UpdateState(ctx, id, state, strings.TrimSpace(failure), now); err != nil {
		return model.Device{}, err
	}
	device.State, device.Failure, device.LastHeartbeat = state, strings.TrimSpace(failure), now
	if err := s.audit(ctx, "device", id, "heartbeat", state); err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (s *Service) ListDevices(ctx context.Context, towerID string) ([]model.Device, error) {
	if strings.TrimSpace(towerID) == "" {
		return nil, model.ErrInvalid
	}
	return s.Devices.ForTower(ctx, towerID)
}

// ReconcileDeviceHealth turns expired heartbeats into durable device alerts.
func (s *Service) ReconcileDeviceHealth(ctx context.Context, towerID string) ([]model.Device, error) {
	devices, err := s.ListDevices(ctx, towerID)
	if err != nil {
		return nil, err
	}
	now := s.Clock.Now()
	for i := range devices {
		if now.Sub(devices[i].LastHeartbeat) <= deviceHeartbeatMaxAge {
			continue
		}
		if err := s.Devices.UpdateState(ctx, devices[i].ID, "stale", "heartbeat overdue", now); err != nil {
			return nil, err
		}
		devices[i].State, devices[i].Failure = "stale", "heartbeat overdue"
		a := alert.Open(now, towerID, "device-heartbeat", model.SeverityWarning, "telemetry gateway "+devices[i].ID+" has not sent a heartbeat")
		if err := s.Alerts.Upsert(ctx, a); err != nil {
			return nil, err
		}
		if err := s.audit(ctx, "device", devices[i].ID, "stale", devices[i].Failure); err != nil {
			return nil, err
		}
	}
	return devices, nil
}
