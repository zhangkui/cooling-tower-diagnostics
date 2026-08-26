package ingest

import (
	"cooling-tower-diagnostics/internal/model"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

type ProtocolDecoder struct{}

func (ProtocolDecoder) DecodeJSON(b []byte) (model.TelemetryFrame, error) {
	var f model.TelemetryFrame
	if err := json.Unmarshal(b, &f); err != nil {
		return f, err
	}
	if f.SentAt.IsZero() {
		f.SentAt = time.Now().UTC()
	}
	return f, nil
}
func (ProtocolDecoder) DecodeBinary(b []byte) (model.TelemetryFrame, error) {
	if len(b) < 16 {
		return model.TelemetryFrame{}, fmt.Errorf("short frame")
	}
	seq := int64(binary.BigEndian.Uint64(b[:8]))
	millis := int64(binary.BigEndian.Uint64(b[8:16]))
	f := model.TelemetryFrame{Sequence: seq, SentAt: time.UnixMilli(millis).UTC(), Fields: map[string]float64{}}
	for i := 16; i+8 <= len(b); i += 8 {
		f.Fields[fmt.Sprintf("channel_%d", (i-16)/8)] = float64(binary.BigEndian.Uint64(b[i:i+8])) / 100
	}
	return f, nil
}
func EncodeBinary(f model.TelemetryFrame) []byte {
	b := make([]byte, 16+len(f.Fields)*8)
	binary.BigEndian.PutUint64(b, uint64(f.Sequence))
	binary.BigEndian.PutUint64(b[8:], uint64(f.SentAt.UnixMilli()))
	i := 16
	for _, v := range f.Fields {
		binary.BigEndian.PutUint64(b[i:], uint64(v*100))
		i += 8
	}
	return b
}
