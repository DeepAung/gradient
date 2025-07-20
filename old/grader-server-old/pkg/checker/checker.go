package checker

import (
	"context"
	"io"
	"os"
	"unicode"
)

const maxReadByte = 512

type CodeChecker interface {
	CheckFile(ctx context.Context, filename1, filename2 string) (bool, error)
	CheckContent(ctx context.Context, reader1, reader2 io.Reader) (bool, error)
}

type codeChecker struct{}

func NewCodeChecker() CodeChecker {
	return codeChecker{}
}

func (c codeChecker) CheckFile(ctx context.Context, filename1, filename2 string) (bool, error) {
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

func (c codeChecker) CheckContent(ctx context.Context, reader1, reader2 io.Reader) (bool, error) {
	_ = ctx

	b1 := make([]byte, maxReadByte)
	idx1, n1, err := c.readFileAndTrimFront(reader1, b1)
	if err != nil && err != io.EOF {
		return false, err
	}

	b2 := make([]byte, maxReadByte)
	idx2, n2, err := c.readFileAndTrimFront(reader2, b2)
	if err != nil && err != io.EOF {
		return false, err
	}

	// in case of EOF, n = 0
	if n1 == 0 || n2 == 0 {
		return (n1 == 0 && n2 == 0), nil
	}

	for {
		breakOuterLoop := false
		for idx1 < n1 && idx2 < n2 {
			if b1[idx1] != b2[idx2] {
				breakOuterLoop = true
				break // break to check for trim back case (e.g. '\n' != ' ')
			}

			idx1++
			idx2++
		}

		if breakOuterLoop {
			break
		}

		if idx1 == n1 {
			idx1 = 0
			n1, err = reader1.Read(b1)
			if err != nil {
				if err == io.EOF {
					break // break to check for trim back case
				}
				return false, err
			}
		}

		if idx2 == n2 {
			idx2 = 0
			n2, err = reader2.Read(b2)
			if err != nil {
				if err == io.EOF {
					break // break to check for trim back case
				}
				return false, err
			}
		}
	}

	// trim back
	if n1 != 0 {
		for idx1 < n1 {
			if unicode.IsSpace(rune(b1[idx1])) {
				idx1++
				continue
			}
			return false, nil
		}
	}
	if n2 != 0 {
		for idx2 < n2 {
			if unicode.IsSpace(rune(b2[idx2])) {
				idx2++
				continue
			}
			return false, nil
		}
	}
	return true, nil
}

func (c codeChecker) readFileAndTrimFront(file io.Reader, b []byte) (idx int, n int, err error) {
	for {
		n, err = file.Read(b)
		if err != nil {
			return
		}

		idx = 0
		for idx < n && unicode.IsSpace(rune(b[idx])) {
			idx++
		}
		if idx != n {
			return
		}
	}
}
