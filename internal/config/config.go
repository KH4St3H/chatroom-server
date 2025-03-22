package config

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ListenAddr            string
	ListenPort            int
	DatabaseDSN           string
	ConnectionWorkerCount int
	ConnectionTimeout     time.Time
	Debug                 bool
}

func setupViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.Set("DEBUG", true)
	viper.Set("LISTEN_ADDR", "0.0.0.0")
	viper.Set("LISTEN_PORT", 15000)
	viper.Set("CONNECTION_TIMEOUT", time.Second*2)
}

func NewConfig() (*Config, error) {
	setupViper()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if !viper.IsSet("database_dsn") {
		return nil, errors.New("database DSN not set")
	}
	return &Config{
		ListenAddr:        viper.GetString("LISTEN_ADDR"),
		ListenPort:        viper.GetInt("LISTEN_PORT"),
		DatabaseDSN:       viper.GetString("DATABASE_DSN"),
		Debug:             viper.GetBool("DEBUG"),
		ConnectionTimeout: viper.GetTime("CONNECTION_TIMEOUT"),
	}, nil
}
