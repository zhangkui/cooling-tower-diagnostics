package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/store"
	"encoding/json"
	"fmt"
	"time"
)

type AuditWriter struct {
	Store store.AuditStore
	Clock func() time.Time
}

func (w AuditWriter) Write(ctx context.Context, entity, id, action string, payload any) error {
	b, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	now := time.Now().UTC()
	if w.Clock != nil {
		now = w.Clock()
	}
	return w.Store.Append(ctx, model.AuditEvent{ID: fmt.Sprintf("audit-%d", now.UnixNano()), Entity: entity, EntityID: id, Action: action, Details: string(b), CreatedAt: now})
}
func AuditPayload(actor, source string) map[string]string {
	return map[string]string{"actor": actor, "source": source}
}
