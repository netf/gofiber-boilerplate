package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
	LogLevel      string
	SentryDSN     string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("LOG_LEVEL", "info")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found: %v", err)
	}

	cfg := &Config{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		LogLevel:      viper.GetString("LOG_LEVEL"),
		SentryDSN:     viper.GetString("SENTRY_DSN"),
	}

	// Validate essential configurations
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	return cfg
}
