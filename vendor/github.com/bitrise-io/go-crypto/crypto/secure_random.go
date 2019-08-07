package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	mathrand "math/rand"
	"time"

	"github.com/pkg/errors"
)

// SecureRandomBytes ...
func SecureRandomBytes(length int64) ([]byte, error) {
	random := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, random); err != nil {
		// this could never could/should happen
		return nil, errors.WithStack(err)
	}
	return random, nil
}

// SecureRandomHex ...
func SecureRandomHex(length int64) (string, error) {
	randomBytes, err := SecureRandomBytes(length)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(string(hex.EncodeToString(randomBytes))), nil
}

// SecureRandomHash ...
func SecureRandomHash(n int) string {
	mathrand.Seed(time.Now().UTC().UnixNano())

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathrand.Intn(len(letterBytes))]
	}
	return string(b)
}
