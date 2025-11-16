package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http       HttpSection       `yaml:"http"`
	Grpc       GrpcSection       `yaml:"grpc"`
	Prometheus PrometheusSection `yaml:"prometheus"`
	Jaeger     JaegerSection     `yaml:"jaeger"`
	Clients    ClientsSection    `yaml:"clients"`
}

type HttpSection struct {
	Address string `yaml:"address"`
}

type GrpcSection struct {
	Address string `yaml:"address"`
}

type PrometheusSection struct {
	Address string `yaml:"address"`
}

type JaegerSection struct {
	Address string `yaml:"address"`
}

type ClientsSection struct {
	Auth struct {
		Address string `yaml:"address"`
	} `yaml:"auth"`
	User struct {
		Address string `yaml:"address"`
	} `yaml:"user"`
}

func GetConfig() *Config {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to parse configuration: %s", err))
	}
	return &config
}
