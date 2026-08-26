package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir  string
	Addr     string
	Database string
}

func Load() Config {
	dir := os.Getenv("DATA_DIR")
	if dir == "" {
		dir = "./data"
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return Config{DataDir: dir, Addr: addr, Database: filepath.Join(dir, "cooling.db")}
}
