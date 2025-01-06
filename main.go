package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal/pokecache"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type pokeAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type config struct {
	nextUrl *string
	prevUrl *string
}

// need to have closure on this variable since it is used later, but not defined yet
var validCommands map[string]cliCommand
var curIndexUrls config
var cache *pokecache.Cache

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	splitStrings := strings.Fields(text)
	return splitStrings
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap() error {
	baseUrl := ""
	if curIndexUrls.nextUrl == nil && curIndexUrls.prevUrl == nil {
		baseUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	if curIndexUrls.nextUrl == nil && curIndexUrls.prevUrl != nil {
		fmt.Println("No more areas to display")
		return nil
	}

	if curIndexUrls.nextUrl != nil {
		baseUrl = *curIndexUrls.nextUrl
	}

	cacheget, hit := cache.Get(baseUrl)
	if !hit {
		//create request
		req, err := http.NewRequest("GET", baseUrl, nil)
		if err != nil {
			return err
		}
		//create header

		//make client and send request
		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		cache.Add(baseUrl, data)
		cacheget = data
	}
	//parse response
	var pokedexAreas pokeAreas
	//decoder := json.NewDecoder(data)
	err := json.Unmarshal(cacheget, &pokedexAreas)
	if err != nil {
		return err
	}

	curIndexUrls.nextUrl = pokedexAreas.Next
	curIndexUrls.prevUrl = pokedexAreas.Previous

	fmt.Println()
	for _, area := range pokedexAreas.Results {
		fmt.Println(area.Name)
	}
	fmt.Println()

	return nil
}

func commandMapb() error {
	baseUrl := ""
	if curIndexUrls.nextUrl == nil && curIndexUrls.prevUrl == nil {
		baseUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	if curIndexUrls.prevUrl == nil && curIndexUrls.nextUrl != nil {
		fmt.Println("No more areas to display")
		return nil
	}

	if curIndexUrls.prevUrl != nil {
		baseUrl = *curIndexUrls.prevUrl
	}

	cacheget, hit := cache.Get(baseUrl)
	if !hit {
		//create request
		req, err := http.NewRequest("GET", baseUrl, nil)
		if err != nil {
			return err
		}
		//create header

		//make client and send request
		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		cache.Add(baseUrl, data)
		cacheget = data
	}
	//parse response
	var pokedexAreas pokeAreas
	//decoder := json.NewDecoder(data)
	err := json.Unmarshal(cacheget, &pokedexAreas)
	if err != nil {
		return err
	}

	curIndexUrls.nextUrl = pokedexAreas.Next
	curIndexUrls.prevUrl = pokedexAreas.Previous

	fmt.Println()
	for _, area := range pokedexAreas.Results {
		fmt.Println(area.Name)
	}
	fmt.Println()

	return nil
}

func commandHelp() error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmdData := range validCommands {
		fmt.Printf("%s: %s\n", cmdData.name, cmdData.description)
	}
	return nil
}

func main() {
	curIndexUrls = config{}
	cache = pokecache.NewCache(5 * time.Minute)

	validCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 areas in the pokedex",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 areas in the pokedex",
			callback:    commandMapb,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		command := input[0]
		//fmt.Printf("Your command was: %s\n", command[0])
		cmdData, exists := validCommands[command]
		if !exists {
			fmt.Printf("Unknown command")
			continue
		}
		//indexUrls := &config{}
		err := cmdData.callback()
		if err != nil {
			fmt.Println("Error executing command: ", err)
		}

		fmt.Println()
	}
}
