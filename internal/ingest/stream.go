package ingest

import (
	"bufio"
	"context"
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"io"
	"sync/atomic"
	"time"
)

type StreamStats struct {
	Frames  atomic.Int64
	Invalid atomic.Int64
	Bytes   atomic.Int64
	Last    time.Time
}

func ConsumeJSONLines(ctx context.Context, r io.Reader, stats *StreamStats, handle func(context.Context, model.ReadingRequest) error) error {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 4096)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		line := scanner.Bytes()
		stats.Bytes.Add(int64(len(line)))
		var req model.ReadingRequest
		if err := json.Unmarshal(line, &req); err != nil {
			stats.Invalid.Add(1)
			continue
		}
		stats.Frames.Add(1)
		stats.Last = time.Now().UTC()
		if err := handle(ctx, req); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func SplitFrames(data []byte, width int) [][]byte {
	if width <= 0 {
		return nil
	}
	out := make([][]byte, 0, (len(data)+width-1)/width)
	for len(data) > 0 {
		n := width
		if n > len(data) {
			n = len(data)
		}
		out = append(out, append([]byte(nil), data[:n]...))
		data = data[n:]
	}
	return out
}
