package checker

import (
	"context"
	"io"
	"os"
	"strings"
)

type simpleCodeChecker struct{}

func NewSimpleCodeChecker() CodeChecker {
	return simpleCodeChecker{}
}

func (c simpleCodeChecker) CheckFile(
	ctx context.Context,
	filename1, filename2 string,
) (bool, error) {
	file1, err := os.Open(filename1)
	if err != nil {
		return false, err
	}
	defer file1.Close()

	file2, err := os.Open(filename2)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	return c.CheckContent(ctx, file1, file2)
}

func (c simpleCodeChecker) CheckContent(
	ctx context.Context,
	reader1, reader2 io.Reader,
) (bool, error) {
	b1, err := io.ReadAll(reader1)
	if err != nil {
		return false, err
	}

	b2, err := io.ReadAll(reader2)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(b1)) == strings.TrimSpace(string(b2)), nil
}
