// Package auth authentication
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ErrInvalidPasswordHash invalid password
var ErrInvalidPasswordHash = errors.New("invalid hash format")

// ErrPasswordMismatch different passwords
var ErrPasswordMismatch = errors.New("incorrect password")

// argon2Params define the algorithm params
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// defaultArgon2Params default params
var defaultArgon2Params = argon2Params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// HashPassword generate a secure hash of the password
func HashPassword(password string) (string, error) {
	salt := make([]byte, defaultArgon2Params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultArgon2Params.iterations,
		defaultArgon2Params.memory,
		defaultArgon2Params.parallelism,
		defaultArgon2Params.keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultArgon2Params.memory,
		defaultArgon2Params.iterations,
		defaultArgon2Params.parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword compare the password in plain text to it hash
func VerifyPassword(password, encodedHash string) error {
	params, salt, hash, err := parseHash(encodedHash)
	if err != nil {
		return fmt.Errorf("parsing the hash: %w", err)
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return ErrPasswordMismatch
	}

	return nil
}

// parseHash extract the params and values of the hash
func parseHash(encodedHash string) (*argon2Params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	params := &argon2Params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism); err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}
	params.saltLength, err = intToUint32(len(salt), "saltLength")
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	params.keyLength, err = intToUint32(len(hash), "hashLength")
	if err != nil {
		return nil, nil, nil, err
	}

	return params, salt, hash, nil
}

func intToUint32(value int, field string) (uint32, error) {
	if value < 0 {
		return 0, fmt.Errorf("%s no puede ser negativo", field)
	}

	if uint64(value) > uint64(math.MaxUint32) {
		return 0, fmt.Errorf(
			"%s excede el máximo permitido para uint32: %d",
			field,
			value,
		)
	}

	return uint32(value), nil
}
