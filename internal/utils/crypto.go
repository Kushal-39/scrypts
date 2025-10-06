package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func GenerateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	return nonce, err
}

func EncryptAESGCM(key, nonce, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}
