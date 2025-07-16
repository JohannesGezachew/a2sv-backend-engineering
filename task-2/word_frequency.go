package main

import (
	"regexp"
	"strings"
)

// WordFrequency takes a string and returns a map with word frequencies
// Words are treated case-insensitively and punctuation is ignored
func WordFrequency(text string) map[string]int {
	if text == "" {
		return make(map[string]int)
	}

	// Convert to lowercase for case-insensitive comparison
	text = strings.ToLower(text)
	
	// Remove punctuation and split by whitespace
	reg := regexp.MustCompile(`[^\w\s]`)
	cleanText := reg.ReplaceAllString(text, "")
	
	// Split into words and filter empty strings
	words := strings.Fields(cleanText)
	
	// Count word frequencies
	frequency := make(map[string]int)
	for _, word := range words {
		if word != "" {
			frequency[word]++
		}
	}
	
	return frequency
}