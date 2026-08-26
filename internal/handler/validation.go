package handler

import (
	"net/http"
	"strconv"
	"strings"
)

type Query struct {
	TowerID string
	Sensor  string
	Minutes int
	State   string
	Limit   int
}

func queryFrom(r *http.Request) Query {
	q := r.URL.Query()
	minutes, _ := strconv.Atoi(q.Get("minutes"))
	if minutes <= 0 {
		minutes = 30
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	return Query{TowerID: strings.TrimSpace(q.Get("tower_id")), Sensor: strings.TrimSpace(q.Get("sensor")), Minutes: minutes, State: strings.TrimSpace(q.Get("state")), Limit: limit}
}

func methodAllowed(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}
func requireJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
func clientWantsCSV(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/csv") || r.URL.Query().Get("format") == "csv"
}
