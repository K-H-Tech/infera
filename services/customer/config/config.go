package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Http       HttpSection       `yaml:"http"`
	Grpc       GrpcSection       `yaml:"grpc"`
	Prometheus PrometheusSection `yaml:"prometheus"`
	Jaeger     JaegerSection     `yaml:"jaeger"`
	Clients    ClientsSection    `yaml:"clients"`
	Postgres   PostgresSection   `yaml:"postgres"`
	GrowthBook GrowthBookSection `yaml:"growthBook"`
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
	Example struct {
		Address string `yaml:"address"`
	} `yaml:"example"`
}

type PostgresSection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type GrowthBookSection struct {
	ClientKey          string        `yaml:"clientKey"`
	APIHost            string        `yaml:"apiHost"`
	PollDatasourceTime time.Duration `yaml:"pollDatasourceTime"`
}

func GetConfig() *Config {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to parse configuration: %s", err))
	}
	return &config
}
