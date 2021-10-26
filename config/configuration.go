package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

type Config struct {
	Hostname string `config:"host"`
	Modem    string `config:"modem"`
	Port     string `config:"port"`
}

func getDefaultConfig() *Config {
	return &Config{
		Hostname: "192.168.2.254",
		Modem:    "TD5130",
		Port:     "9618",
	}
}

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
