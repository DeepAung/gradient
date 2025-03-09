package evaluator

import (
	"io"
	"strings"
	"unicode"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(reader1, reader2 io.Reader) (bool, error) {
	b1, err := io.ReadAll(reader1)
	if err != nil {
		return false, err
	}

	b2, err := io.ReadAll(reader2)
	if err != nil {
		return false, err
	}

	lines1 := strings.Split(strings.TrimSpace(string(b1)), "\n")
	lines2 := strings.Split(strings.TrimSpace(string(b2)), "\n")

	for i := 0; ; i++ {
		if i >= len(lines1) && i >= len(lines2) {
			return true, nil
		} else if i >= len(lines1) {
			if e.isAllSpaces(lines2[i]) {
				continue
			}
			return false, nil
		} else if i >= len(lines2) {
			if e.isAllSpaces(lines1[i]) {
				continue
			}
			return false, nil
		} else {
			line1, line2 := strings.TrimSpace(lines1[i]), strings.TrimSpace(lines2[i])
			if line1 == line2 {
				continue
			}
			return false, nil
		}
	}
}

func (e *Evaluator) isAllSpaces(str string) bool {
	for _, char := range str {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}
