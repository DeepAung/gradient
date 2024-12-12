package utils

import (
	"crypto/aes"
	"fmt"

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

func Encrypt(str string, secret []byte) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	encryptedBytes := make([]byte, len(str))
	block.Encrypt(encryptedBytes, []byte(str))

	return string(encryptedBytes), nil
}

func Decrypt(str string, secret []byte) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	decryptedBytes := make([]byte, len(str))
	block.Decrypt(decryptedBytes, []byte(str))

	return string(decryptedBytes), nil
}
