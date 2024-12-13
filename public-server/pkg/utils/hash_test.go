package utils_test

import (
	"testing"
	"testing/quick"

	"github.com/DeepAung/gradient/public-server/pkg/utils"
)

func TestEncryptDecrypt(t *testing.T) {
	assertion := func(input string) bool {
		original := input
		secret := []byte("abcdefghijklmnop")

		encrypted, err := utils.Encrypt(original, secret)
		if err != nil {
			return false
		}

		decrypted, err := utils.Decrypt(encrypted, secret)
		if err != nil {
			return false
		}

		return original == decrypted
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Fatal(err)
	}
}
