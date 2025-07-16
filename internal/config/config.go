package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP HTTPConfig
	DB DBConfig
	Log LogConfig
}

type HTTPConfig struct{
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeSeconds int `mapstructure:"conn_max_lifetime_seconds"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("APP")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.BindEnv("db.dsn")
	_ = v.BindEnv("http.port")
	_ = v.BindEnv("log.level")

	v.SetDefault("http.port", "8080")
	v.SetDefault("db.max_open_conns", 10)
	v.SetDefault("db.max_idle_conns", 5)
	v.SetDefault("db.conn_max_lifetime_seconds", 300)
	v.SetDefault("log.level", "info")

	v.SetConfigName("config")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err !=  nil {
		return nil, err
	}

	return &cfg, nil
}