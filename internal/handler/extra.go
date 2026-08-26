package handler

import (
	"cooling-tower-diagnostics/internal/model"
	"encoding/json"
	"net/http"
	"strconv"
)

func decodeInt(v string, def int) int {
	n, e := strconv.Atoi(v)
	if e != nil || n <= 0 {
		return def
	}
	return n
}
func writeCSV(w http.ResponseWriter, data []byte, name string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	w.WriteHeader(200)
	_, _ = w.Write(data)
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"code": code, "message": message})
}
func stateLabel(s string) string {
	switch s {
	case model.AlertOpen:
		return "开放"
	case model.AlertAcknowledged:
		return "已确认"
	case model.AlertEscalated:
		return "已升级"
	case model.AlertClosed:
		return "已关闭"
	}
	return s
}
