package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "removes punctuation and handles spaces",
			input:    " Hello, World! ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "handles no spaces with punctuation",
			input:    "Hello, World!",
			expected: []string{"hello", "world"},
		},
		{
			name:     "handles all lowercase with punctuation",
			input:    "hello, world!",
			expected: []string{"hello", "world"},
		},
		{
			name:     "handles single word",
			input:    "  exit  ",
			expected: []string{"exit"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := cleanInput(c.input)
			if len(actual) != len(c.expected) {
				t.Errorf("Expected %v, but got %v", c.expected, actual)
			}
			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				if word != expectedWord {
					t.Errorf("Expected %v, but got %v", expectedWord, word)
				}
			}
		})
	}
}
