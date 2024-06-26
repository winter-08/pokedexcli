package pokedexapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

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

  url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
  res, err := http.Get(url)

  if err != nil {
    fmt.Printf("There was an error getting location area: %v\n", err)
    return nil, err
  }


  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    return nil, errors.New(fmt.Sprintf("Status: %v\nMessage: %s\n", res.StatusCode, string(body)))
  }

  cache.Add(name, body)
  
  response := LocationArea{}
  json.Unmarshal(body, &response)

  return &response, nil

}

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
    return nil, errors.New(fmt.Sprintf("Status: %v\nMessage: %s\n", res.StatusCode, string(body)))
  }

  cache.Add(url, body)

  response := Locations{}
  json.Unmarshal(body, &response)

  return &response, nil

}

type Pokemon struct {
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
  Stats []struct {
    BaseStat int `json:"base_stat"`
    Effort int `json:"effort"`
    Stat struct {
      Name string `json:"name"`
      URL string `json:"url"`
    } `json:"stat"`
  } `json:"stats"`
  Types []struct {
    Slot int `json:"slot"`
    Type struct {
      Name string `json:"name"`
      URL string `json:"url"`
    } `json:"type"`
  } `json:"types"`
}

func GetPokemon(name string, cache *pokecache.Cache) (*Pokemon, error) {
  cachedRes, found := cache.Get(name)
  if found == true {
    cachedResponse := Pokemon{}
    json.Unmarshal(*cachedRes, &cachedResponse)
    return &cachedResponse, nil
  }

  url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
  res, err := http.Get(url)

  if err != nil {
    fmt.Printf("There was an error getting pokemon: %v\n", err)
    return nil, err
  }


  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    return nil, errors.New(fmt.Sprintf("Status: %v\nMessage: %s\n", res.StatusCode, string(body)))
  }

  cache.Add(name, body)
  
  response := Pokemon{}
  json.Unmarshal(body, &response)

  return &response, nil

}
