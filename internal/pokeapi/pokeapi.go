package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type pokeMap struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url string) (pokeMap, error) {
	res, err := http.Get(url)
	if err != nil {
		return pokeMap{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return pokeMap{}, err
	}
	params := pokeMap{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return pokeMap{}, err
	}
	return params, nil
}
