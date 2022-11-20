package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Log  Log  `yaml:"log"`
	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`
}

type Log struct {
	Level string `yaml:"level"`
	Mode  string `yaml:"mode"`
}

// New reads application configuration from specified file path.
func New(path string) (*Settings, error) {
	var config Settings

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
