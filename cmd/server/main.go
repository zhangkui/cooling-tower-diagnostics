package main

import (
	"context"
	"cooling-tower-diagnostics/internal/clock"
	"cooling-tower-diagnostics/internal/config"
	"cooling-tower-diagnostics/internal/diagnostics"
	"cooling-tower-diagnostics/internal/handler"
	"cooling-tower-diagnostics/internal/ingest"
	"cooling-tower-diagnostics/internal/service"
	"cooling-tower-diagnostics/internal/store"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	db, err := store.Open(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	if err := db.EnsureOptionalSchema(ctx); err != nil {
		log.Fatal(err)
	}
	c := clock.Real{}
	t := store.TowerStore{DB: db}
	r := store.ReadingStore{DB: db}
	a := store.AlertStore{DB: db}
	au := store.AuditStore{DB: db}
	th := store.ThresholdStore{DB: db}
	sensors := store.SensorStore{DB: db}
	maintenance := store.MaintenanceStore{DB: db}
	queries := store.QueryStore{DB: db}
	devices := store.DeviceStore{DB: db}
	e := diagnostics.Engine{Towers: t, Readings: r, Thresholds: th, Clock: c}
	s := &service.Service{Towers: t, Readings: r, Alerts: a, Audits: au, Thresholds: th, Engine: e, Clock: c, Decoder: ingest.Decoder{}, Sensors: sensors, Maintenance: maintenance, Queries: queries, Devices: devices}
	if err := s.Seed(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("cooling tower diagnostics listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler.New(s).Routes()); err != nil {
		log.Fatal(err)
	}
}
