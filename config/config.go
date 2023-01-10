package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Listen string `json:"listen" ini:"listen" cfg:"listen" cfgDefault:"0.0.0.0:2211"`
}

func Load() (Config, error) {
	var cfg = Config{}
	config.PrefixEnv = "MOOCA"
	config.File = "config.ini"
	err := config.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
