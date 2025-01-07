package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal/pokecache"
	"io"
	"math/rand"
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

type pokePokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
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
var myPokemon map[string]pokePokemon

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

func commandCatch(params ...string) error {
	if strings.Join(params, "") == "" {
		return fmt.Errorf("explore command requires an area id or name")
	}

	if len(params) > 1 {
		return fmt.Errorf("explore command only takes one parameter")
	}

	baseUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", params[0])

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
	var pokemon pokePokemon
	//decoder := json.NewDecoder(data)
	err := json.Unmarshal(cacheget, &pokemon)
	if err != nil {
		return err
	}

	pokemonBExp := pokemon.BaseExperience
	pokemonName := pokemon.Name

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	time.Sleep(2 * time.Second)
	catchChance := 80 - (pokemonBExp / 1000)
	if catchChance < 10 {
		catchChance = 10
	}
	catchChint := int(catchChance)
	randNum := rand.Intn(100)
	if randNum < catchChint {
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	myPokemon[pokemonName] = pokemon

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

func commandInspect(params ...string) error {
	if strings.Join(params, "") == "" {
		return fmt.Errorf("inspect command requires a pokemon id or name")
	}

	if len(params) > 1 {
		return fmt.Errorf("inspect command only takes one parameter")
	}

	if _, exists := myPokemon[params[0]]; !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}

	hp := 0
	attack := 0
	defense := 0
	specialAttack := 0
	specialDefense := 0
	speed := 0

	for _, val := range myPokemon[params[0]].Stats {
		switch val.Stat.Name {
		case "hp":
			hp = val.BaseStat
		case "attack":
			attack = val.BaseStat
		case "defense":
			defense = val.BaseStat
		case "special-attack":
			specialAttack = val.BaseStat
		case "special-defense":
			specialDefense = val.BaseStat
		case "speed":
			speed = val.BaseStat
		}
	}
	inspectedPokemon := myPokemon[params[0]]
	fmt.Printf("Name: %s\n", inspectedPokemon.Name)
	fmt.Printf("Height: %d\n", inspectedPokemon.Height)
	fmt.Printf("Weight: %d\n", inspectedPokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf("\t-hp: %d\n", hp)
	fmt.Printf("\t-attack: %d\n", attack)
	fmt.Printf("\t-defense: %d\n", defense)
	fmt.Printf("\t-special-attack: %d\n", specialAttack)
	fmt.Printf("\t-special-defense: %d\n", specialDefense)
	fmt.Printf("\t-speed: %d\n", speed)
	fmt.Println("Types:")
	for _, val := range inspectedPokemon.Types {
		fmt.Printf("\t-%s\n", val.Type.Name)
	}

	return nil
}

func commandPokedex(params ...string) error {
	if strings.Join(params, "") != "" {
		return fmt.Errorf("pokedex command does not take any parameters")
	}
	fmt.Println()
	fmt.Println("Your Pokedex:")
	for idx := range myPokemon {
		fmt.Printf(" - %v\n", idx)
	}
	fmt.Println()

	return nil
}

func main() {
	curIndexUrls = config{}
	myPokemon = map[string]pokePokemon{}
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
		"catch": {
			name:        "catch <pokemon id or name>",
			description: "Attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon id or name>",
			description: "Displays the stats of a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all caught pokemon",
			callback:    commandPokedex,
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
