package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Port        int    `json:"port" ini:"port" cfg:"port" cfgDefault:"2211"`
	DatabaseURL string `json:"database_url" ini:"database_url" cfg:"database_url" cfgDefault:"postgres://postgres:postgres@localhost:5432/mooca?sslmode=disable"`
}

func Load() (Config, error) {
	var cfg = Config{}
	config.File = "config.ini"
	err := config.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
