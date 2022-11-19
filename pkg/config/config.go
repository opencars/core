package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Log  Log  `yaml:"log"`
	GRPC GRPC `yaml:"grpc"`
}

type GRPC struct {
	Registrations ServiceGRPC `yaml:"registrations"`
	Operations    ServiceGRPC `yaml:"operations"`
	VinDecoding   ServiceGRPC `yaml:"vin_decoding"`
	ALPR          ServiceGRPC `yaml:"alpr"`
}

type ServiceGRPC struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s *ServiceGRPC) Address() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// Log represents settings for application logger.
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
