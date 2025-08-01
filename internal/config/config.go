package config

import (
	"encoding/json"
	"os"
)

type URLConfig struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Data   string `json:"data"`
}

type Config struct {
	URLs        []string    `json:"urls"`
	Interval    int         `json:"interval"`
	Timeout     int         `json:"timeout"`
	Logfile     string      `json:"logfile"`
	Verbose     bool        `json:"verbose"`
	Export      string      `json:"export"`
	Concurrency int         `json:"concurrency"`
	URLConfigs  []URLConfig `json:"url_configs"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 10
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 5
	}
	if cfg.Interval == 0 {
		cfg.Interval = 0
	}
	return &cfg, nil
}
