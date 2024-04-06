package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"gssm/types"
	"math/big"
)

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSecretKey() (string, string) {
	key := make([]rune, 40)

	mx := big.NewInt(int64(len(alphabet)))

	for i := 0; i < len(key); i++ {
		b, _ := rand.Int(rand.Reader, mx)
		key[i] = rune(alphabet[b.Int64()])
	}

	secret := string(key)

	for i := 2; i < len(key); i++ {
		key[i] = '*'
	}

	mask := string(key)

	return secret, mask
}

func GenerateAccessKey() string {
	key := make([]rune, 20)

	mx := big.NewInt(int64(len(alphabet)))

	for i := 0; i < len(key); i++ {
		b, _ := rand.Int(rand.Reader, mx)
		key[i] = rune(alphabet[b.Int64()])
	}

	secret := string(key)

	return "GSSM" + secret
}

func GetSignature(userId, projectId types.ULID, secretKey string) ([]byte, error) {
	data := []byte("user-" + userId + "-project-" + projectId)

	h := hmac.New(sha512.New, []byte(secretKey))

	_, err := h.Write(data)

	if err != nil {
		return nil, err
	}

	dataHmac := h.Sum(nil)

	return dataHmac, nil
}
