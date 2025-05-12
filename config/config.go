package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Configs *Config

type MYSQL struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
}

type Config struct {
	MYSQL MYSQL `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

func init() {
	fmt.Println("svc config initializing...")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic("Config Read failed: " + err.Error())
	}
	err := viper.Unmarshal(&Configs)
	if err != nil {
		panic("Config decode failed: " + err.Error())
	}
}
