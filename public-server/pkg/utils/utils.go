package utils

import (
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

func Compare(str, hashedStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil
}
