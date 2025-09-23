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
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous *string
}

var commands map[string]cliCommand
var location config
var cache *pokecache.Cache

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
	}

	location = config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: nil,
	}

	cache = pokecache.NewCache(time.Second * 7)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			cmd := cleanInput(scanner.Text())[0]
			val, ok := commands[cmd]
			if ok {
				err := val.callback(&location)
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

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *config) error {
	params := pokeapi.PokeMap{}
	val, ok := cache.Get(c.Next)
	if ok {
		err := json.Unmarshal(val, &params)
		if err != nil {
			return err
		}
	} else {
		p, err := pokeapi.GetLocationAreas(c.Next)
		if err != nil {
			return err
		}
		params = p
		jsonData, _ := json.Marshal(params)
		cache.Add(c.Next, jsonData)
	}
	fmt.Println(params.Next)
	c.Next = params.Next
	c.Previous = params.Previous
	for _, r := range params.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func commandMapb(c *config) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	params := pokeapi.PokeMap{}
	val, ok := cache.Get(*c.Previous)
	if ok {
		err := json.Unmarshal(val, &params)
		if err != nil {
			return err
		}
	} else {
		p, err := pokeapi.GetLocationAreas(*c.Previous)
		if err != nil {
			return err
		}
		params = p
		jsonData, _ := json.Marshal(params)
		cache.Add(*c.Previous, jsonData)
	}
	c.Next = params.Next
	c.Previous = params.Previous
	for _, r := range params.Results {
		fmt.Println(r.Name)
	}
	return nil
}
