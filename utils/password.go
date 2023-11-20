package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"golang.org/x/crypto/pbkdf2"
)

const (
	minSaltSize    = 8
	iterationCount = 210000
	keyLength      = 64
)

type HashResult struct {
	Hash []byte
	Salt []byte
}

func generateSalt() []byte {
	saltBytes := make([]byte, minSaltSize)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}

	return saltBytes
}

func HashPassword(password string) HashResult {
	saltString := generateSalt()
	salt := bytes.NewBuffer(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, iterationCount, keyLength, sha512.New)
	return HashResult{Hash: df, Salt: saltString}
}

func VerifyPassword(password string, hash []byte, salt []byte) bool {
	saltBytes := bytes.NewBuffer(salt).Bytes()
	df := pbkdf2.Key([]byte(password), saltBytes, iterationCount, keyLength, sha512.New)

	return equal(hash, df)
}

func equal(oldHash, newHash []byte) bool {
	diff := uint64(len(oldHash)) ^ uint64(len(newHash))

	for i := 0; i < len(oldHash) && i < len(newHash); i++ {
		diff |= uint64(oldHash[i]) ^ uint64(newHash[i])
	}

	return diff == 0
}
