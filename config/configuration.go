package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type Config struct {
	Host  string `config:"host"`
	Modem string `config:"modem"`
	Port  string `config:"port"`
	User  string `config:"user"`
	Pass  string `config:"pass"`
}

func Load() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		file.NewOptionalBackend("/etc/modem_exporter.yaml"),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := Config{
		Host:  "192.168.2.254",
		Modem: "TD5130",
		Port:  "9618",
		User:  "",
		Pass:  "",
	}

	err := loader.Load(context.Background(), &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
