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
	Redis      RedisSection      `yaml:"redis"`
	KavehNegar KavehNegarSection `yaml:"kavehnegar"`
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

type RedisSection struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KavehNegarSection struct {
	ApiKey   string `yaml:"apikey"`
	Url      string `yaml:"url"`
	Template string `yaml:"template"`
}

func GetConfig() *Config {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to parse configuration: %s", err))
	}
	return &config
}
