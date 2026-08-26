package alert

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"time"
)

func Open(now time.Time, tower, rule, severity, message string) model.Alert {
	return model.Alert{ID: fmt.Sprintf("alert-%d", now.UnixNano()), TowerID: tower, Rule: rule, Severity: severity, State: model.AlertOpen, Message: message, OpenedAt: now, UpdatedAt: now}
}
func Acknowledge(a *model.Alert, now time.Time) error {
	if a.State != model.AlertOpen && a.State != model.AlertEscalated {
		return model.ErrConflict
	}
	a.State = model.AlertAcknowledged
	a.UpdatedAt = now
	a.AcknowledgedAt = &now
	return nil
}
func Escalate(a *model.Alert, now time.Time) error {
	if a.State != model.AlertOpen {
		return model.ErrConflict
	}
	a.State = model.AlertEscalated
	a.UpdatedAt = now
	return nil
}
func Close(a *model.Alert, now time.Time) error {
	if a.State == model.AlertClosed {
		return model.ErrConflict
	}
	a.State = model.AlertClosed
	a.UpdatedAt = now
	a.ClosedAt = &now
	return nil
}
