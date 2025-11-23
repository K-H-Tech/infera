package configuration

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Http       HttpSection       `yaml:"http"`
	Grpc       GrpcSection       `yaml:"grpc"`
	Prometheus PrometheusSection `yaml:"prometheus"`
	Jaeger     JaegerSection     `yaml:"jaeger"`
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

func LoadConfig(serviceName string) *Config {
	//viper := viper.New()

	// Load from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Current directory
	viper.AddConfigPath("./services/" + serviceName + "/")
	viper.AddConfigPath("./configuration") // Optional configuration directory

	// Environment variable support
	//viper.SetEnvPrefix(strings.ToUpper(serviceName))       // Environment variables must be prefixed with SERVICE_NAME
	viper.AutomaticEnv()                                   // Automatically override with ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // db.host -> DB_HOST

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No local configuration.yaml found. Relying on environment variables.")
		}
	} else {
		log.Printf("Using configuration file: %s", viper.ConfigFileUsed())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to parse configuration: %s", err))
	}

	return &config
}
