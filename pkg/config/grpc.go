package config

import "strconv"

type ServiceGRPC struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPC struct {
	Registrations ServiceGRPC `yaml:"registrations"`
	Operations    ServiceGRPC `yaml:"operations"`
	VinDecoding   ServiceGRPC `yaml:"vin_decoding"`
}

func (s *ServiceGRPC) Address() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}
