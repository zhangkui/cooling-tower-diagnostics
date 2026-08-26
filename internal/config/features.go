package config

import "os"

type Features struct {
	Seed          bool
	EnableReplay  bool
	EnableExports bool
	WorkerCount   int
}

func LoadFeatures() Features {
	f := Features{Seed: true, EnableReplay: true, EnableExports: true, WorkerCount: 4}
	if os.Getenv("NO_SEED") == "1" {
		f.Seed = false
	}
	if os.Getenv("NO_REPLAY") == "1" {
		f.EnableReplay = false
	}
	return f
}
