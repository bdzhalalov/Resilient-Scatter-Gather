package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Addr     string `mapstructure:"ADDR"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

var Cfg Config

func InitConfig() Config {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}

	return Cfg
}
