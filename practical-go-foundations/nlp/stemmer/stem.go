package stemmer

// Stem takes a word and returns its stemmed version.
func Stem(word string) string {
	// A simple implementation of a stemmer
	// This is a placeholder for an actual stemming algorithm
	if len(word) < 3 {
		return word
	}

	// Remove common suffixes
	suffixes := []string{"ing", "ed", "ly", "es", "s"}
	for _, suffix := range suffixes {
		if len(word) > len(suffix) && word[len(word)-len(suffix):] == suffix {
			return word[:len(word)-len(suffix)]
		}
	}

	return word
}
