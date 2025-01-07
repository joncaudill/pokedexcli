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
	callback    func(params ...string) error
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

type pokeLocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	Id        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int `json:"chance"`
				ConditionValues []struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"condition_values"`
				MaxLevel int `json:"max_level"`
				Method   struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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

func commandExit(params ...string) error {
	if strings.Join(params, "") != "" {
		return fmt.Errorf("exit command does not take any parameters")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandExplore(params ...string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	if strings.Join(params, "") == "" {
		return fmt.Errorf("explore command requires an area id or name")
	}

	if len(params) > 1 {
		return fmt.Errorf("explore command only takes one parameter")
	}

	baseUrl = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", params[0])

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
	var pokedexLocationAreas pokeLocationArea
	//decoder := json.NewDecoder(data)
	err := json.Unmarshal(cacheget, &pokedexLocationAreas)
	if err != nil {
		return err
	}

	fmt.Println()
	for _, encounter := range pokedexLocationAreas.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	fmt.Println()

	return nil
}

func commandMap(params ...string) error {
	if strings.Join(params, "") != "" {
		return fmt.Errorf("map command does not take any parameters")
	}
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

func commandMapb(params ...string) error {
	if strings.Join(params, "") != "" {
		return fmt.Errorf("mapb command does not take any parameters")
	}
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

func commandHelp(params ...string) error {
	if strings.Join(params, "") != "" {
		return fmt.Errorf("help command does not take any parameters")
	}
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
		"explore": {
			name:        "explore <area id or name>",
			description: "Displays the pokemon in a given area",
			callback:    commandExplore,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		param1 := ""
		if len(input) == 0 {
			continue
		}
		if len(input) > 1 {
			param1 = input[1]
		}
		command := input[0]
		//fmt.Printf("Your command was: %s\n", command[0])
		cmdData, exists := validCommands[command]
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}
		//indexUrls := &config{}
		err := cmdData.callback(param1)
		if err != nil {
			fmt.Println("Error executing command: ", err)
		}

		fmt.Println()
	}
}
