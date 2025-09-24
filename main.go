package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal/pokeapi"
	"internal/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	Next     string
	Previous *string
}

var commands map[string]cliCommand
var location config
var cache *pokecache.Cache
var baseURL string

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Displays the names of the previous 20 location areas in the Pokemon world if there are any",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Allows user to list all pokemon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Allows user to try to catch a pokemon",
			callback:    commandCatch,
		},
	}

	baseURL = "https://pokeapi.co/api/v2/location-area/"
	location = config{
		Next:     baseURL + "?offset=0&limit=20",
		Previous: nil,
	}

	cache = pokecache.NewCache(time.Second * 7)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			params := cleanInput(scanner.Text())
			cmd := params[0]
			param_1 := ""
			if len(params) >= 2 {
				param_1 = params[1]
			}
			val, ok := commands[cmd]
			if ok {
				err := val.callback(&location, param_1)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config, input string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, input string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *config, input string) error {
	return getAndDisplayLocationAreas(c, c.Next)
}

func commandMapb(c *config, input string) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return getAndDisplayLocationAreas(c, *c.Previous)
}

func getAndDisplayLocationAreas(c *config, url string) error {
	params := pokeapi.PokeMap{}
	val, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(val, &params)
		if err != nil {
			return err
		}
	} else {
		p, err := pokeapi.Get[pokeapi.PokeMap](url)
		if err != nil {
			return err
		}
		params = p
		jsonData, _ := json.Marshal(params)
		cache.Add(url, jsonData)
	}
	c.Next = params.Next
	c.Previous = params.Previous
	for _, r := range params.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func commandExplore(c *config, input string) error {
	if input == "" {
		return nil
	}
	url := baseURL + input
	fmt.Printf("Exploring %s...\n", input)
	params := pokeapi.LocationAreaPokemon{}
	val, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(val, &params)
		if err != nil {
			return err
		}
	} else {
		p, err := pokeapi.Get[pokeapi.LocationAreaPokemon](url)
		if err != nil {
			fmt.Println("Problem with finding location area. Please make sure that it is spelled correctly.")
			return err
		}
		params = p
		jsonData, _ := json.Marshal(params)
		cache.Add(url, jsonData)
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range params.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, input string) error {
	if input == "" {
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", input)
	return nil
}
