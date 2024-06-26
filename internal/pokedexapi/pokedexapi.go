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

type LocationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationArea(name string, cache *pokecache.Cache) (*LocationArea, error) {
  cachedRes, found := cache.Get(name)

  if found == true {
    cachedResponse := LocationArea{}
    json.Unmarshal(*cachedRes, &cachedResponse)
    return &cachedResponse, nil
  }

  url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", name)
  fmt.Printf("url: %v\n", url)
  res, err := http.Get(url)

  if err != nil {
    fmt.Printf("There was an error getting location area: %v\n", err)
    return nil, err
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
    fmt.Printf("body: %v\n", body)
    return nil, errors.New("Error getting location area")
  }

  cache.Add(name, body)
  
  response := LocationArea{}
  json.Unmarshal(body, &response)

  return &response, nil

}

func GetLocations(url string, cache *pokecache.Cache) (*Locations, error) {
  cachedRes, found := cache.Get(url)

  fmt.Printf("found: %v\n", found)

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

  cache.Add(url, body)

  response := Locations{}
  json.Unmarshal(body, &response)

  return &response, nil

}
