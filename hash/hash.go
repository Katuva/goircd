package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"strings"
)

const (
	time    = 1
	memory  = 1024
	threads = 4
	keyLen  = 32
	saltLen = 16
)

func Make(plaintext string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plaintext), salt, time, memory, threads, keyLen)

	encoded := base64.RawStdEncoding.EncodeToString(salt) + "$" +
		base64.RawStdEncoding.EncodeToString(hash)

	return encoded, nil
}

func Verify(plaintext string, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	testHash := argon2.IDKey([]byte(plaintext), salt, time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(hash, testHash) == 1
}
