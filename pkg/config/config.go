package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Log  Log  `yaml:"log"`
	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`
	NATS NATS `yaml:"nats"`
}

// NATS contains configuration details for application event API.
type NATS struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Address returns calculated address for connecting to NATS.
func (nats *NATS) Address() string {
	if nats.User != "" && nats.Password != "" {
		return fmt.Sprintf("nats://%s:%s@%s:%d", nats.User, nats.Password, nats.Host, nats.Port)
	}

	return fmt.Sprintf("nats://%s:%d", nats.Host, nats.Port)
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
