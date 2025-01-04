package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	splitStrings := strings.Fields(text)
	for word := range splitStrings {
		splitStrings[word] = strings.ToLower(splitStrings[word])
	}
	return splitStrings
}

func main() {
	fmt.Println("Hello, World!")
}
