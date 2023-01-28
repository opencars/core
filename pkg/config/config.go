package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Log  Log  `yaml:"log"`
	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`
	NATS NATS `yaml:"nats"`
}

type NodeNATS struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (node *NodeNATS) Address(user, password string) string {
	if user != "" && password != "" {
		return fmt.Sprintf("nats://%s:%s@%s:%d", user, password, node.Host, node.Port)
	}

	return fmt.Sprintf("nats://%s:%d", node.Host, node.Port)
}

// NATS contains configuration details for application event API.
type NATS struct {
	Nodes    []NodeNATS `yaml:"nodes"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
}

// Address returns calculated address for connecting to NATS.
func (nats *NATS) Address() string {
	addrs := make([]string, 0)

	for _, node := range nats.Nodes {
		addrs = append(addrs, node.Address(nats.User, nats.Password))

	}

	return strings.Join(addrs, ",")
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
