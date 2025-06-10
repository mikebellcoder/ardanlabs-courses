package main

import (
	"bufio"
	"fmt"
	"maps"
	"slices"
	"sort"

	"os"
	"regexp"
	"strings"
)

var wordRe = regexp.MustCompile(`\w+`) // use [a-zA-Z]+ if /w+ does not work properly

func main() {
	file, err := os.Open("sherlock.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	freq := make(map[string]int) // word -> count
	s := bufio.NewScanner(file)

	for s.Scan() {
		words := wordRe.FindAllString(s.Text(), -1)
		for _, word := range words {
			freq[strings.ToLower(word)]++
		}
	}
	if err := s.Err(); err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	top := topN(freq, 10)
	fmt.Println("Top 10 words:", top)
}

func topN(freq map[string]int, n int) []string {
	words := slices.Collect(maps.Keys(freq))

	sort.Slice(words, func(i, j int) bool {
		// Sort in desc order
		return freq[words[i]] > freq[words[j]]
	})

	n = min(n, len(words))
	return words[:n]
}

func mapDemo() {
	heroes := map[string]string{
		"Superman":     "Clark",
		"Batman":       "Bruce",
		"Wonder Woman": "Diana",
	}

	// Keys
	for k := range heroes {
		fmt.Println(k)
	}

	// Keys + Values
	for k, v := range heroes {
		fmt.Println(v, "is", k)
	}
}
