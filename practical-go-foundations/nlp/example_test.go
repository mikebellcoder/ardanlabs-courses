package nlp_test

import (
	"fmt"

	"github.com/mikebellcoder/nlp"
)

func ExampleTokenize() {
	// Example usage of the Tokenize function
	text := "Who's on first?"
	tokens := nlp.Tokenize(text)
	for _, token := range tokens {
		fmt.Println(token)
	}
	// Output:
	// who
	// s
	// on
	// first
}
