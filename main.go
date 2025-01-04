package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	splitStrings := strings.Fields(text)
	return splitStrings
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		fmt.Printf("Your command was: %s\n", input[0])
	}
}
