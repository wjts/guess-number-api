package argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2ID struct {
	format  string
	version int
	time    uint32
	memory  uint32
	keyLen  uint32
	saltLen uint32
	threads uint8
}

func DefaultParams() Argon2ID {
	return Argon2ID{
		format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version: argon2.Version,
		time:    1,
		memory:  65536,
		keyLen:  32,
		saltLen: 16,
		threads: 4,
	}
}

func (argon Argon2ID) Hash(plain string) (string, error) {
	salt := make([]byte, argon.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plain), salt, argon.time, argon.memory, argon.threads, argon.keyLen)

	return fmt.Sprintf(
			argon.format,
			argon.version,
			argon.memory,
			argon.time,
			argon.threads,
			base64.RawStdEncoding.EncodeToString(salt),
			base64.RawStdEncoding.EncodeToString(hash),
		),
		nil
}

func (argon Argon2ID) Verify(plain, hash string) (bool, error) {
	hashParts := strings.Split(hash, "$")

	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &argon.memory, &argon.time, &argon.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, argon.time, argon.memory, argon.threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1, nil
}
