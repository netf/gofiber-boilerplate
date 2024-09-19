package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/dgrijalva/jwt-go"
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
	PrivateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}

	PwSalt = make([]byte, 32)
	_, err = rand.Read(PwSalt)
	return err
}
