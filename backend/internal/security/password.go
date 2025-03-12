package security

import (
    "crypto/rand"
    "encoding/base64"
    "golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
}

var defaultConfig = &PasswordConfig{
    time:    1,
    memory:  64 * 1024,
    threads: 4,
    keyLen:  32,
}

func HashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, 
        defaultConfig.time,
        defaultConfig.memory,
        defaultConfig.threads,
        defaultConfig.keyLen)

    encodedHash := base64.RawStdEncoding.EncodeToString(hash)
    encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

    return encodedSalt + ":" + encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
    parts := strings.Split(encodedHash, ":")
    if len(parts) != 2 {
        return false, fmt.Errorf("invalid hash format")
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[0])
    if err != nil {
        return false, err
    }

    hash := argon2.IDKey([]byte(password), salt,
        defaultConfig.time,
        defaultConfig.memory,
        defaultConfig.threads,
        defaultConfig.keyLen)

    encodedNewHash := base64.RawStdEncoding.EncodeToString(hash)
    return encodedNewHash == parts[1], nil
}