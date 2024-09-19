package config

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/netf/gofiber-boilerplate/internal/errors"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/spf13/viper"
)

var (
	PrivateKey *ecdsa.PrivateKey
	PwSalt     []byte
)

func InitAuth() error {
	var err error
	PrivateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}

	PwSalt = make([]byte, 32)
	_, err = rand.Read(PwSalt)
	return err
}

type Config struct {
	ServerAddress  string
	DatabaseURL    string
	LogLevel       string
	SentryDSN      string
	Environment    string
	AuthPrivateKey string
	AuthSalt       string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENVIRONMENT", "dev")

	// Try to read the config file, but don't return an error if it's not found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// If the config file is not found, log a warning and continue
		log.Warn().Msg("Config file not found. Using environment variables and defaults.")
	}

	cfg := &Config{
		ServerAddress:  viper.GetString("SERVER_ADDRESS"),
		DatabaseURL:    viper.GetString("DATABASE_URL"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		SentryDSN:      viper.GetString("SENTRY_DSN"),
		Environment:    viper.GetString("ENVIRONMENT"),
		AuthPrivateKey: viper.GetString("AUTH_PRIVATE_KEY"),
		AuthSalt:       viper.GetString("AUTH_SALT"),
	}

	// Validate essential configurations
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required but not set")
	}

	// We don't need to validate AuthPrivateKey and AuthSalt here
	// as they will be generated if not provided

	return cfg, nil
}
