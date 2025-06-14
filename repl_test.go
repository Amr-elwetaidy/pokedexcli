package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: " Hello, World! ", expected: []string{"hello,", "world!"}},
		{input: "Hello, World!", expected: []string{"hello,", "world!"}},
		{input: "hello, world!", expected: []string{"hello,", "world!"}},
	}

	for _, c := range cases {
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
	}
}
