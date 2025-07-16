package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "single character",
			input:    "a",
			expected: true,
		},
		{
			name:     "simple palindrome",
			input:    "racecar",
			expected: true,
		},
		{
			name:     "not palindrome",
			input:    "hello",
			expected: false,
		},
		{
			name:     "palindrome with spaces",
			input:    "race car",
			expected: true,
		},
		{
			name:     "palindrome with punctuation",
			input:    "A man, a plan, a canal: Panama",
			expected: true,
		},
		{
			name:     "palindrome with mixed case",
			input:    "Madam",
			expected: true,
		},
		{
			name:     "complex palindrome",
			input:    "Was it a car or a cat I saw?",
			expected: true,
		},
		{
			name:     "not palindrome with punctuation",
			input:    "Hello, world!",
			expected: false,
		},
		{
			name:     "numeric palindrome",
			input:    "12321",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}