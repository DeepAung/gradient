package evaluator

import (
	"strings"
	"testing"

	"github.com/DeepAung/gradient/website-server/pkg/asserts"
)

func TestEvaluate(t *testing.T) {
	eval := NewEvaluator()

	tests := []struct {
		name     string
		str1     string
		str2     string
		expected bool
	}{
		{
			name:     "normal, equal",
			str1:     "a\nb b",
			str2:     "a\nb b",
			expected: true,
		},
		{
			name:     "end of line suffix, equal",
			str1:     "a\nb b\n\n\n\n",
			str2:     "a\nb b",
			expected: true,
		},
		{
			name:     "inline suffix, equal",
			str1:     "a\nb b     \nc",
			str2:     "a\nb b\nc",
			expected: true,
		},
		{
			name:     "multiple new lines, not equal",
			str1:     "a\nb",
			str2:     "a\n\n\nb",
			expected: false,
		},
		{
			name:     "empty reader1, not equal",
			str1:     "",
			str2:     "abc",
			expected: false,
		},
		{
			name:     "empty reader2, not equal",
			str1:     "abc",
			str2:     "",
			expected: false,
		},
		{
			name:     "empty reader1 and reader2, equal",
			str1:     "",
			str2:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		reader1 := strings.NewReader(tt.str1)
		reader2 := strings.NewReader(tt.str2)
		got, err := eval.Evaluate(reader1, reader2)

		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "eval", got, tt.expected)
	}
}
