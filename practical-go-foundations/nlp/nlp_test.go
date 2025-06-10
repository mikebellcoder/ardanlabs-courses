package nlp

import (
	"slices"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Who's on first?", []string{"who", "s", "on", "first"}},
		{"Hello, World!", []string{"hello", "world"}},
		{"Go is awesome.", []string{"go", "is", "awesome"}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := Tokenize(test.input)
			if !slices.Equal(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
