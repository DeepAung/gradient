package checker

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"
)

func TestCheckContent(t *testing.T) {
	tests := []struct {
		name          string
		reader1       io.Reader
		reader2       io.Reader
		expectedCheck bool
		expectedError error
	}{
		{
			name:          "no trim, equal",
			reader1:       strings.NewReader("a\nb\nccc\nd   e"),
			reader2:       strings.NewReader("a\nb\nccc\nd   e"),
			expectedCheck: true,
			expectedError: nil,
		},
		{
			name:          "no trim, not equal",
			reader1:       strings.NewReader("a\nb\nccc\nd   e"),
			reader2:       strings.NewReader("a\nb\nccc\nd      e"),
			expectedCheck: false,
			expectedError: nil,
		},
		{
			name:          "trim front, equal",
			reader1:       strings.NewReader("\n\n   \nabc\ndefg"),
			reader2:       strings.NewReader("abc\ndefg"),
			expectedCheck: true,
			expectedError: nil,
		},
		{
			name:          "trim front, not equal",
			reader1:       strings.NewReader(" abc\ndefg"),
			reader2:       strings.NewReader("aabc\ndefg"),
			expectedCheck: false,
			expectedError: nil,
		},
		{
			name:          "trim front and back, equal",
			reader1:       strings.NewReader("    abc\ndefg    \n\n\n\n  "),
			reader2:       strings.NewReader("\nabc\ndefg\n"),
			expectedCheck: true,
			expectedError: nil,
		},
		{
			name:          "trim front and back, not equal",
			reader1:       strings.NewReader("    abc\ndefg    \n\n\n\n  "),
			reader2:       strings.NewReader("abc\n\ndefg"),
			expectedCheck: false,
			expectedError: nil,
		},
		{
			name:          "empty reader1, not equal",
			reader1:       strings.NewReader(""),
			reader2:       strings.NewReader("aaaa"),
			expectedCheck: false,
			expectedError: nil,
		},
		{
			name:          "empty reader2, not equal",
			reader1:       strings.NewReader("aaaa"),
			reader2:       strings.NewReader(""),
			expectedCheck: false,
			expectedError: nil,
		},
		{
			name:          "empty reader1 and reader2, equal",
			reader1:       strings.NewReader(""),
			reader2:       strings.NewReader(""),
			expectedCheck: true,
			expectedError: nil,
		},
		{
			name: "more than maxReadByte, equal",
			reader1: func() io.Reader {
				s1 := make([]byte, 30000)
				for i := 0; i < 300; i++ {
					s1[i] = ' '
				}
				for i := 300; i < 30000; i++ {
					s1[i] = 'a'
				}
				return bytes.NewReader(s1)
			}(),
			reader2: func() io.Reader {
				s2 := make([]byte, 30000)
				for i := 0; i < 30000-300; i++ {
					s2[i] = 'a'
				}
				for i := 30000 - 300; i < 30000; i++ {
					s2[i] = ' '
				}
				return bytes.NewReader(s2)
			}(),
			expectedCheck: true,
			expectedError: nil,
		},
	}

	checker := NewCodeChecker()
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check, err := checker.CheckContent(ctx, tt.reader1, tt.reader2)

			if check != tt.expectedCheck {
				t.Fatalf("invalid check, expect=%v, got=%v", tt.expectedCheck, check)
			}
			if err != tt.expectedError {
				t.Fatalf("invalid error, expect=%v, got=%v", tt.expectedError, err)
			}
		})
	}
}

func TestCheckContentEqual(t *testing.T) {
	assertion := func(val int) bool {
		str := generateRandomSpace(1024) +
			generateRandomString(rand.Intn(100000)) +
			generateRandomSpace(1024)
		reader1 := strings.NewReader(str)
		reader2 := strings.NewReader(str)
		ctx := context.Background()

		checker := NewCodeChecker()
		check, err := checker.CheckContent(ctx, reader1, reader2)
		if err != nil {
			return false
		}

		simpleChecker := NewSimpleCodeChecker()
		simpleCheck, err := simpleChecker.CheckContent(ctx, reader1, reader2)
		if err != nil {
			return false
		}

		return check == simpleCheck
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Fatal("failed")
	}
}

func TestCheckContentNotEqual(t *testing.T) {
	assertion := func(val int) bool {
		str1 := generateRandomSpace(rand.Intn(1024)) +
			generateRandomString(rand.Intn(100000)) +
			generateRandomSpace(rand.Intn(1024))
		str2 := generateRandomSpace(rand.Intn(1024)) +
			generateRandomString(rand.Intn(100000)) +
			generateRandomSpace(rand.Intn(1024))
		reader1 := strings.NewReader(str1)
		reader2 := strings.NewReader(str2)
		ctx := context.Background()

		checker := NewCodeChecker()
		check, err := checker.CheckContent(ctx, reader1, reader2)
		if err != nil {
			return false
		}

		simpleChecker := NewSimpleCodeChecker()
		simpleCheck, err := simpleChecker.CheckContent(ctx, reader1, reader2)
		if err != nil {
			return false
		}

		return check == simpleCheck
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Fatal("failed")
	}
}

var (
	letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890 \n")
	spaceBytes  = []byte("\t\n\v\f\r ") // from unicode.IsSpace()
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateRandomSpace(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = spaceBytes[rand.Intn(len(spaceBytes))]
	}
	return string(b)
}
