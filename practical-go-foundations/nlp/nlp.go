package nlp

import (
	"regexp"
	"strings"
)

var (
	// "Who's on first?" -> [Who s on first]
	wordRe = regexp.MustCompile(`\w+`)
)

// Tokenize splits the input text into words and returns them as a slice of strings. All words will be in lowercase.
func Tokenize(text string) []string {
	// Find all words in the text
	words := wordRe.FindAllString(text, -1)
	var tokens []string
	for _, word := range words {
		// Convert to lowercase
		token := strings.ToLower(word)
		tokens = append(tokens, token)
	}
	return tokens
}
