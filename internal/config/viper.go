package config

import (
	"github.com/milobella/enviper"
	"github.com/spf13/viper"
)

func ReadConfiguration() (*Configuration, error) {
	e := enviper.New(viper.New())
	e.SetEnvPrefix("ABILITY")

	e.SetConfigName("config")
	e.SetConfigType("toml")
	e.AddConfigPath("/etc/ability/")
	e.AddConfigPath("$HOME/.ability")
	e.AddConfigPath(".")
	if err := e.ReadInConfig(); err != nil {
		return nil, err
	}

	var C Configuration
	//_ = e.Unmarshal(&C)
	if err := e.Unmarshal(&C); err != nil {
		return nil, err
	}
	return &C, nil
}
