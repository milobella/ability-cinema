package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Configuration
type Configuration struct {
	Server ServerConfiguration          `mapstructure:"server"`
	Tools  map[string]ToolConfiguration `mapstructure:"tools"`
}

// fun String() : Serialization function of Configuration (for logging)
func (c Configuration) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Configuration serialization error %s", err)
	}
	return string(b)
}

type ServerConfiguration struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type ToolConfiguration struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
