package utils

import (
	"github.com/netf/gofiber-boilerplate/config"
	"golang.org/x/crypto/scrypt"
)

func PwHash(pw []byte) ([]byte, error) {
	return scrypt.Key(pw, config.PwSalt, 32768, 8, 1, 32)
}
