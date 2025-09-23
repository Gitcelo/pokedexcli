package pokeapi

import (
	"encoding/json"
	"io"
	"log"
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

func GetLocationAreas(url string) (PokeMap, error) {
	res, err := http.Get(url)
	if err != nil {
		return PokeMap{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return PokeMap{}, err
	}
	params := PokeMap{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return PokeMap{}, err
	}
	return params, nil
}
