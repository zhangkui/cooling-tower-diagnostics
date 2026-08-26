package handler

import (
	"cooling-tower-diagnostics/internal/model"
	"cooling-tower-diagnostics/internal/service"
	"cooling-tower-diagnostics/internal/store"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HTTP struct{ S *service.Service }

func New(s *service.Service) *HTTP { return &HTTP{S: s} }
func (h *HTTP) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", h.health)
	m.HandleFunc("/", h.index)
	m.HandleFunc("/api/towers", h.towers)
	m.HandleFunc("/api/readings", h.readings)
	m.HandleFunc("/api/replay", h.replay)
	m.HandleFunc("/api/thresholds", h.thresholds)
	m.HandleFunc("/api/diagnose", h.diagnose)
	m.HandleFunc("/api/alerts", h.alerts)
	m.HandleFunc("/api/audit", h.audit)
	m.HandleFunc("/api/sensors", h.sensors)
	m.HandleFunc("/api/maintenance", h.maintenance)
	m.HandleFunc("/api/stats", h.stats)
	m.HandleFunc("/api/report", h.report)
	m.HandleFunc("/api/forecast", h.forecast)
	m.HandleFunc("/api/capacity", h.capacity)
	m.HandleFunc("/api/water-balance", h.waterBalance)
	m.HandleFunc("/api/devices", h.devices)
	m.HandleFunc("/api/devices/", h.devices)
	return m
}
func (h *HTTP) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (h *HTTP) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	template.Must(template.New("index").Parse(page)).Execute(w, nil)
}
func (h *HTTP) towers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		v, e := h.S.ListTowers(r.Context())
		respond(w, v, e)
	case http.MethodPost:
		var q model.CreateTowerRequest
		if decode(r, &q) != nil {
			respondErr(w, model.ErrInvalid)
			return
		}
		v, e := h.S.CreateTower(r.Context(), q)
		respond(w, v, e)
	default:
		http.Error(w, "method not allowed", 405)
	}
}
func (h *HTTP) readings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var q model.ReadingRequest
	if decode(r, &q) != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	v, e := h.S.IngestReading(r.Context(), q)
	respond(w, v, e)
}
func (h *HTTP) replay(w http.ResponseWriter, r *http.Request) {
	var q model.ReplayRequest
	if decode(r, &q) != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	respond(w, map[string]string{"status": "accepted"}, h.S.Replay(r.Context(), q))
}
func (h *HTTP) thresholds(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.S.ListThresholds(r.Context())
		respond(w, v, e)
		return
	}
	var q model.ThresholdRequest
	if decode(r, &q) != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	respond(w, map[string]string{"status": "saved"}, h.S.ConfigureThreshold(r.Context(), q))
}
func (h *HTTP) diagnose(w http.ResponseWriter, r *http.Request) {
	var q model.DiagnosisRequest
	q.TowerID = r.URL.Query().Get("tower_id")
	q.SinceMinutes = 30
	if v := r.URL.Query().Get("minutes"); v != "" {
		for _, x := range []int{1, 5, 15, 30, 60, 120} {
			if v == string(rune('0'+x)) {
				q.SinceMinutes = x
			}
		}
	}
	v, e := h.S.EvaluateAndAlert(r.Context(), q)
	respond(w, v, e)
}
func (h *HTTP) alerts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.S.ListAlerts(r.Context(), r.URL.Query().Get("state"))
		respond(w, v, e)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		respondErr(w, model.ErrInvalid)
		return
	}
	var q model.StateRequest
	_ = decode(r, &q)
	respond(w, map[string]string{"status": "updated"}, h.S.TransitionAlert(r.Context(), parts[2], q.Note))
}
func (h *HTTP) audit(w http.ResponseWriter, r *http.Request) {
	v, e := h.S.ListAudit(r.Context(), r.URL.Query().Get("entity"))
	respond(w, v, e)
}

func (h *HTTP) sensors(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.S.ListSensors(r.Context(), r.URL.Query().Get("tower_id"))
		respond(w, v, e)
		return
	}
	var q model.CreateSensorRequest
	if decode(r, &q) != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	v, e := h.S.SaveSensor(r.Context(), q)
	respond(w, v, e)
}

func (h *HTTP) maintenance(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.S.ListMaintenance(r.Context(), r.URL.Query().Get("tower_id"))
		respond(w, v, e)
		return
	}
	if r.Method == http.MethodPatch {
		respond(w, map[string]string{"status": "completed"}, h.S.CompleteMaintenance(r.Context(), r.URL.Query().Get("id")))
		return
	}
	var q model.CreateMaintenanceRequest
	if decode(r, &q) != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	v, e := h.S.PlanMaintenance(r.Context(), q)
	respond(w, v, e)
}

func (h *HTTP) stats(w http.ResponseWriter, r *http.Request) {
	v, e := h.S.TowerStats(r.Context(), r.URL.Query().Get("tower_id"))
	respond(w, v, e)
}

func (h *HTTP) report(w http.ResponseWriter, r *http.Request) {
	q := queryFrom(r)
	now := time.Now().UTC()
	filter := model.ReportFilter{TowerID: q.TowerID, From: now.Add(-time.Duration(q.Minutes) * time.Minute), To: now}
	if q.Sensor != "" {
		filter.Sensors = []string{q.Sensor}
	}
	if clientWantsCSV(r) {
		data, result, err := (&service.Exporter{Queries: store.QueryStore{DB: h.S.Readings.DB}, Engine: h.S.Engine}).CSV(r.Context(), filter)
		if err != nil {
			respondErr(w, err)
			return
		}
		writeCSV(w, data, result.Filename)
		return
	}
	v, e := h.S.Report(r.Context(), filter)
	respond(w, v, e)
}

func (h *HTTP) forecast(w http.ResponseWriter, r *http.Request) {
	q := queryFrom(r)
	v, e := h.S.Forecast(r.Context(), q.TowerID, q.Sensor, decodeInt(r.URL.Query().Get("horizon"), 12))
	respond(w, v, e)
}

func (h *HTTP) capacity(w http.ResponseWriter, r *http.Request) {
	q := queryFrom(r)
	rated, _ := strconv.ParseFloat(r.URL.Query().Get("rated_flow"), 64)
	v, e := h.S.Capacity(r.Context(), q.TowerID, rated)
	respond(w, v, e)
}
func (h *HTTP) waterBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	from, fromErr := time.Parse(time.RFC3339, query.Get("from"))
	to, toErr := time.Parse(time.RFC3339, query.Get("to"))
	if fromErr != nil || toErr != nil {
		respondErr(w, model.ErrInvalid)
		return
	}
	tolerance, err := strconv.ParseFloat(query.Get("tolerance_m3"), 64)
	if err != nil && query.Get("tolerance_m3") != "" {
		respondErr(w, model.ErrInvalid)
		return
	}
	v, e := h.S.AssessWaterBalance(r.Context(), model.WaterBalanceRequest{TowerID: query.Get("tower_id"), From: from, To: to, Tolerance: tolerance, MaxGapMinutes: decodeInt(query.Get("max_gap_minutes"), 15)})
	respond(w, v, e)
}
func (h *HTTP) devices(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.S.ListDevices(r.Context(), r.URL.Query().Get("tower_id"))
		respond(w, v, e)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if r.Method == http.MethodPost && len(parts) == 2 {
		var q model.CreateDeviceRequest
		if decode(r, &q) != nil {
			respondErr(w, model.ErrInvalid)
			return
		}
		v, e := h.S.RegisterDevice(r.Context(), q)
		respond(w, v, e)
		return
	}
	if r.Method == http.MethodPost && len(parts) == 4 && parts[3] == "heartbeat" {
		var q model.HeartbeatRequest
		_ = decode(r, &q)
		v, e := h.S.HeartbeatDevice(r.Context(), parts[2], q.Failure)
		respond(w, v, e)
		return
	}
	if r.Method == http.MethodPost && len(parts) == 4 && parts[3] == "reconcile" {
		v, e := h.S.ReconcileDeviceHealth(r.Context(), parts[2])
		respond(w, v, e)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func respond(w http.ResponseWriter, v any, e error) {
	if e != nil {
		respondErr(w, e)
		return
	}
	writeJSON(w, 200, v)
}
func respondErr(w http.ResponseWriter, e error) {
	code := 500
	if errors.Is(e, model.ErrInvalid) {
		code = 400
	}
	if errors.Is(e, model.ErrNotFound) {
		code = 404
	}
	if errors.Is(e, model.ErrConflict) {
		code = 409
	}
	writeJSON(w, code, map[string]string{"error": e.Error()})
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

const page = `<!doctype html><html><head><meta charset="utf-8"><title>Cooling Tower Diagnostics</title><style>body{font-family:system-ui;margin:40px;background:#f4f6f8;color:#17202a}main{max-width:960px;margin:auto}section{background:white;border:1px solid #d8dee4;padding:20px;margin:14px 0}button{padding:8px 14px;background:#1769aa;color:#fff;border:0}pre{background:#111;color:#b9f6ca;padding:12px;min-height:80px}</style></head><body><main><h1>Cooling Tower Diagnostics</h1><section><button onclick="loadTowers()">Refresh towers</button><button onclick="loadAlerts()">Refresh alerts</button><pre id="out">Ready</pre></section></main><script>async function loadTowers(){out.textContent=JSON.stringify(await (await fetch('/api/towers')).json(),null,2)}async function loadAlerts(){out.textContent=JSON.stringify(await (await fetch('/api/alerts')).json(),null,2)}</script></body></html>`
