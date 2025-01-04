package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	//construct cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello wOrld ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "chariZARd  BULBAsaur ",
			expected: []string{"charizard", "bulbasaur"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}
	//run tests
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %d elements, got %d", len(c.expected), len(actual))
			t.Fail()
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Expected %s, got %s", c.expected[i], actual[i])
				t.Fail()
			}
		}
	}
}
