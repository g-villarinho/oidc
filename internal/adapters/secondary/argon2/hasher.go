package argon2

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"golang.org/x/crypto/argon2"
)

type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

type hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewHasher() ports.Hasher {
	return &hasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,         // 3 iterações
		parallelism: 2,         // 2 threads
		saltLength:  16,        // 16 bytes de salt
		keyLength:   32,        // 32 bytes de chave
	}
}

func (h *hasher) Hash(ctx context.Context, plaintext string) (string, error) {
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.memory,
		h.iterations,
		h.parallelism,
		saltB64,
		hashB64,
	)

	return encodedHash, nil
}

func (h *hasher) Compare(ctx context.Context, plaintext string, hashed string) error {
	params, salt, hash, err := decodeHash(hashed)
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	compareHash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	if subtle.ConstantTimeCompare(hash, compareHash) != 1 {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func decodeHash(encodedHash string) (*hashParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("incompatible argon2 variant")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid version: %w", err)
	}

	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible argon2 version")
	}

	params := &hashParams{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid hash: %w", err)
	}

	params.keyLength = uint32(len(hash))

	return params, salt, hash, nil
}
