package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const cost = 10

func Hash(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), cost)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %v", err)
	}
	return string(hashed), nil
}

func CompareHash(str, hashedStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil
}

func Encrypt(inputStr string, secret []byte) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	gcmInstance, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	encryptedBytes := gcmInstance.Seal(nonce, nonce, []byte(inputStr), nil)
	encrypedStr := base64.StdEncoding.EncodeToString(encryptedBytes)
	return encrypedStr, nil
}

func Decrypt(encrypedStr string, secret []byte) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypedStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	gcmInstance, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]

	decryptedBytes, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}
