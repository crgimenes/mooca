package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Port int `json:"port" ini:"port" cfg:"port" cfgDefault:"2211"`
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
