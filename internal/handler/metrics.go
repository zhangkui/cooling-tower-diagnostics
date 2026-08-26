package handler

import (
	"net/http"
	"strconv"
	"strings"
)

type MetricsView interface{ Snapshot() map[string]any }

func MetricsHandler(metrics MetricsView) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", 405)
			return
		}
		snapshot := metrics.Snapshot()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		for key, value := range snapshot {
			switch v := value.(type) {
			case int64:
				w.Write([]byte(key + " " + strconv.FormatInt(v, 10) + "\n"))
			case int:
				w.Write([]byte(key + " " + strconv.Itoa(v) + "\n"))
			default:
				w.Write([]byte(key + " " + strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(toString(v)), " ", "_")) + "\n"))
			}
		}
	}
}
func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	default:
		return "0"
	}
}
