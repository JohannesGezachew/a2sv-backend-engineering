package main

import "fmt"

func main() {
	// Example usage of WordFrequency
	text := "Hello world! This is a test. Hello again, world."
	freq := WordFrequency(text)
	fmt.Println("Word frequencies:")
	for word, count := range freq {
		fmt.Printf("%s: %d\n", word, count)
	}
	
	fmt.Println()
	
	// Example usage of IsPalindrome
	testStrings := []string{
		"racecar",
		"A man, a plan, a canal: Panama",
		"hello world",
		"Was it a car or a cat I saw?",
	}
	
	fmt.Println("Palindrome checks:")
	for _, str := range testStrings {
		result := IsPalindrome(str)
		fmt.Printf("'%s' is palindrome: %t\n", str, result)
	}
}