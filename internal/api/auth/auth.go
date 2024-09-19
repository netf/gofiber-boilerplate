package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

var (
	PrivateKey *ecdsa.PrivateKey
	PwSalt     []byte
)

type Claims struct {
	jwt.StandardClaims
	Special string `json:"spc,omitempty"`
}

func InitAuth() error {
	var err error

	// Load or generate private key
	PrivateKey, err = loadOrGeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to load or generate private key: %w", err)
	}

	// Load or generate salt
	PwSalt, err = loadOrGenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to load or generate salt: %w", err)
	}

	return nil
}

func loadOrGeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	// Check if the private key is provided as an environment variable
	privateKeyPEM := os.Getenv("AUTH_PRIVATE_KEY")
	if privateKeyPEM != "" {
		block, _ := pem.Decode([]byte(privateKeyPEM))
		if block == nil {
			return nil, fmt.Errorf("failed to decode PEM block containing private key")
		}
		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		log.Info().Msg("Loaded private key from environment variable")
		return privateKey, nil
	}

	// Check if the private key file path is provided
	privateKeyPath := os.Getenv("AUTH_PRIVATE_KEY_FILE")
	if privateKeyPath != "" {
		keyBytes, err := ioutil.ReadFile(privateKeyPath)
		if err == nil {
			block, _ := pem.Decode(keyBytes)
			if block == nil {
				return nil, fmt.Errorf("failed to decode PEM block containing private key")
			}
			privateKey, err := x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
			log.Info().Msg("Loaded private key from file")
			return privateKey, nil
		}
	}

	// If no key is provided, generate a new one
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	log.Info().Msg("Generated new private key")
	return privateKey, nil
}

func loadOrGenerateSalt() ([]byte, error) {
	// Check if the salt is provided as an environment variable
	saltBase64 := os.Getenv("AUTH_SALT")
	if saltBase64 != "" {
		salt, err := base64.StdEncoding.DecodeString(saltBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode salt from environment variable: %w", err)
		}
		if len(salt) != 32 {
			return nil, fmt.Errorf("invalid salt length: expected 32 bytes, got %d", len(salt))
		}
		log.Info().Msg("Loaded salt from environment variable")
		return salt, nil
	}

	// Generate new salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	log.Info().Msg("Generated new salt")
	return salt, nil
}
