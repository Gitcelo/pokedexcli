package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeMap struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaPokemon struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
}

func Get[T any](url string) (T, error) {
	var result T
	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return result, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return result, err
	}
	params := result
	err = json.Unmarshal(body, &params)
	if err != nil {
		return result, err
	}
	return params, nil
}
