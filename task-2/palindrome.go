package main

import (
	"regexp"
	"strings"
)

// IsPalindrome checks if a string is a palindrome
// Ignores spaces, punctuation, and capitalization
func IsPalindrome(text string) bool {
	if text == "" {
		return true
	}

	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Remove all non-alphanumeric characters
	reg := regexp.MustCompile(`[^a-z0-9]`)
	cleanText := reg.ReplaceAllString(text, "")
	
	// Check if the cleaned string reads the same forwards and backwards
	length := len(cleanText)
	for i := 0; i < length/2; i++ {
		if cleanText[i] != cleanText[length-1-i] {
			return false
		}
	}
	
	return true
}