package config

import "strconv"

type HTTP struct {
	Statisfy ServiceHTTP `yaml:"statisfy"`
}

type ServiceHTTP struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Secret string `yaml:"secret"`
	Token  string `yaml:"token"`
}

func (s *ServiceHTTP) Address() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}
