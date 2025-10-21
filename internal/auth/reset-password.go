package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func GenerateSecureToken(nBytes int) (token string, raw []byte, err error) {
	raw = make([]byte, nBytes)
	if _, err = rand.Read(raw); err != nil {
		return "", nil, err
	}

	token = base64.RawURLEncoding.EncodeToString(raw)
	return token, raw, nil
}

func HashToken(raw []byte) []byte {
	sum := sha256.Sum256(raw)
	return sum[:]
}

var (
	ErrWeakPassword = errors.New("senha fraca")
)
