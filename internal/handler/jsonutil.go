package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func readJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return io.EOF
	}
	return json.Unmarshal(body, dst)
}
func writeJSONIndent(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(value)
}
func acceptJSON(r *http.Request) bool {
	v := r.Header.Get("Accept")
	return v == "" || strings.Contains(v, "application/json")
}
func requestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	return "local"
}
