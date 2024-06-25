package pokedexapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Region struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"region"`
		Names []struct {
			Name     string `json:"name"`
			Language struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"language"`
		} `json:"names"`
		GameIndices []struct {
			GameIndex  int `json:"game_index"`
			Generation struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"generation"`
		} `json:"game_indices"`
		Areas []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"areas"`
	} `json:"results"`
}

func GetLocations(url string, cache *pokecache.Cache) (*Locations, error) {
  cachedRes, found := cache.Get(url)

  if found == true {
    cachedResponse := Locations{}
    json.Unmarshal(*cachedRes, &cachedResponse)
    return &cachedResponse, nil
  }

  res, err := http.Get(url)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return nil, err
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
    fmt.Printf("body: %v\n", body)
    return nil, errors.New("Error getting locations")
  }

  response := Locations{}
  json.Unmarshal(body, &response)

  return &response, nil

}
