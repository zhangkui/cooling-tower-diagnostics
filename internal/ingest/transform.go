package ingest

import (
	"cooling-tower-diagnostics/internal/model"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
)

type Transformer struct {
	UnitMap map[string]string
	Alias   map[string]string
}

func NewTransformer() *Transformer {
	return &Transformer{UnitMap: map[string]string{"celsius": "C", "degree_c": "C", "microsiemens": "uS/cm"}, Alias: map[string]string{"cond": "conductivity", "temperature_out": "outlet_temp", "temperature_in": "inlet_temp"}}
}
func (t *Transformer) Sensor(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if v, ok := t.Alias[name]; ok {
		return v
	}
	return name
}
func (t *Transformer) Unit(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if v, ok := t.UnitMap[name]; ok {
		return v
	}
	return name
}
func (t *Transformer) Reading(r model.ReadingRequest) model.ReadingRequest {
	r.Sensor = t.Sensor(r.Sensor)
	r.Unit = t.Unit(r.Unit)
	return r
}
func (t *Transformer) Frame(f model.TelemetryFrame) []model.ReadingRequest {
	out := []model.ReadingRequest{}
	for name, value := range f.Fields {
		out = append(out, model.ReadingRequest{TowerID: f.DeviceID, Sensor: t.Sensor(name), Value: value, Unit: "", RecordedAt: &f.SentAt})
	}
	return out
}
func ParseHex(payload string) ([]byte, error) {
	payload = strings.TrimSpace(payload)
	if strings.HasPrefix(payload, "0x") {
		payload = payload[2:]
	}
	if len(payload)%2 != 0 {
		return nil, fmt.Errorf("odd hex payload")
	}
	return hex.DecodeString(payload)
}
func DecodeFloat32(data []byte) float64 {
	if len(data) < 4 {
		return math.NaN()
	}
	bits := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	return float64(math.Float32frombits(bits))
}
func EncodeFloat32(value float64) []byte {
	bits := math.Float32bits(float32(value))
	return []byte{byte(bits >> 24), byte(bits >> 16), byte(bits >> 8), byte(bits)}
}
func Age(at, timeNow time.Time) time.Duration {
	if timeNow.Before(at) {
		return 0
	}
	return timeNow.Sub(at)
}
