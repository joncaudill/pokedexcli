package main

import (
	"fmt"
	"testing"
	"time"

	"internal/pokecache"
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

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
