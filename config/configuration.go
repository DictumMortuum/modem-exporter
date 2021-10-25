package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

// Config is the exporter CLI configuration.
type Config struct {
	Hostname string        `config:"host"`
	Port     string        `config:"port"`
	Interval time.Duration `config:"interval"`
}

func getDefaultConfig() *Config {
	return &Config{
		Hostname: "192.168.2.254",
		Port:     "6102",
		Interval: 180 * time.Second,
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
