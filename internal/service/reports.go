package service

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ReportCache struct {
	items   map[string][]byte
	expires map[string]time.Time
}

func NewReportCache() *ReportCache {
	return &ReportCache{items: map[string][]byte{}, expires: map[string]time.Time{}}
}
func (c *ReportCache) Get(key string, now time.Time) ([]byte, bool) {
	data, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if now.After(c.expires[key]) {
		delete(c.items, key)
		delete(c.expires, key)
		return nil, false
	}
	return append([]byte(nil), data...), true
}
func (c *ReportCache) Put(key string, data []byte, ttl time.Duration, now time.Time) {
	c.items[key] = append([]byte(nil), data...)
	c.expires[key] = now.Add(ttl)
}
func (c *ReportCache) Clear() { c.items = map[string][]byte{}; c.expires = map[string]time.Time{} }

func ReportKey(filter model.ReportFilter) string {
	return fmt.Sprintf("%s|%s|%s|%s", filter.TowerID, filter.From.Format(time.RFC3339), filter.To.Format(time.RFC3339), strings.Join(filter.Sensors, ","))
}
func MarshalReport(r model.DiagnosticReport) ([]byte, error) { return json.Marshal(r) }
func BuildFilter(tower string, minutes int, sensors []string, now time.Time) model.ReportFilter {
	if minutes <= 0 {
		minutes = 30
	}
	return model.ReportFilter{TowerID: tower, From: now.Add(-time.Duration(minutes) * time.Minute), To: now, Sensors: sensors}
}
func RunReport(ctx context.Context, build func(context.Context, model.ReportFilter) (model.DiagnosticReport, error), filter model.ReportFilter) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	report, err := build(ctx, filter)
	if err != nil {
		return nil, err
	}
	return MarshalReport(report)
}
