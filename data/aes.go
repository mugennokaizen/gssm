package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"io"
)

type AesProcessor struct {
	MasterKey []byte
}

func NewAesProcessor(_ *do.Injector) (*AesProcessor, error) {
	key := []byte(viper.GetString("jwt.refresh_cookie_name"))

	if len(key) != 32 {
		panic(errors.New("master AES key must be length of 32"))
	}
	return &AesProcessor{
		MasterKey: key,
	}, nil
}

func (a *AesProcessor) Encrypt(value string) (string, error) {

	text := []byte(value)

	c, err := aes.NewCipher(a.MasterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(gcm.Seal(nonce, nonce, text, nil)), err
}

func (a *AesProcessor) Decrypt(s string) (string, error) {

	data, err := hex.DecodeString(s)

	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(a.MasterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("data length less than noncesize")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
