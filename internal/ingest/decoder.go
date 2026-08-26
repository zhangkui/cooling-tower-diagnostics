package ingest

import (
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"fmt"
	"strings"
)

type Decoder struct{}

func (Decoder) Decode(payload []byte) (model.ReadingRequest, error) {
	var r model.ReadingRequest
	if err := json.Unmarshal(payload, &r); err != nil {
		return r, fmt.Errorf("decode telemetry: %w", err)
	}
	r.Sensor = strings.ToLower(strings.TrimSpace(r.Sensor))
	if r.TowerID == "" || r.Sensor == "" {
		return r, model.ErrInvalid
	}
	return r, nil
}
func (Decoder) DecodeBatch(payload []byte) ([]model.ReadingRequest, error) {
	var rs []model.ReadingRequest
	if err := json.Unmarshal(payload, &rs); err != nil {
		return nil, err
	}
	for i := range rs {
		rs[i].Sensor = strings.ToLower(strings.TrimSpace(rs[i].Sensor))
		if rs[i].TowerID == "" || rs[i].Sensor == "" {
			return nil, model.ErrInvalid
		}
	}
	return rs, nil
}
