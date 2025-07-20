package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed config_template.yaml
var configTemplate []byte

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
}

type App struct {
	Address string `yaml:"address"`
}

type Database struct {
	DSN string `yaml:"dsn"`
}

func NewConfig() *Config {
	var cfg Config
	if err := yaml.Unmarshal(configTemplate, &cfg); err != nil {
		log.Fatalf("Failed to parsing config yaml file: %v", err)
	}
	return &cfg
}
