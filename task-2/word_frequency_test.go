package main

import (
	"reflect"
	"testing"
)

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:     "single word",
			input:    "hello",
			expected: map[string]int{"hello": 1},
		},
		{
			name:     "multiple words",
			input:    "hello world hello",
			expected: map[string]int{"hello": 2, "world": 1},
		},
		{
			name:     "case insensitive",
			input:    "Hello HELLO hello",
			expected: map[string]int{"hello": 3},
		},
		{
			name:     "with punctuation",
			input:    "Hello, world! Hello world.",
			expected: map[string]int{"hello": 2, "world": 2},
		},
		{
			name:     "complex sentence",
			input:    "The quick brown fox jumps over the lazy dog. The dog was lazy!",
			expected: map[string]int{"the": 3, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 2, "dog": 2, "was": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WordFrequency(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("WordFrequency(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}