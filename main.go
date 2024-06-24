package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type config struct {
  next string 
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(*config) error
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

func main() {

  config := config {
    next: "",
    previous: "",
  }

  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := scanner.Text()


    commands := cliCommands()

    if cmd, ok := commands[input]; ok {
      cmd.callback(&config)
    }

  }

}

func cliCommands() map[string]cliCommand {
  return map[string]cliCommand{
    "help": {
      name: "help",
  description: "Displays a help message",
  callback: commandHelp,
    },
    "exit": {
  name: "exit",
  description: "Exit the pokedex",
  callback: commandExit,
    },
    "map": {
  name: "map",
  description: "Returns 20 locations, calling it subsequently will return the next 20 locations",
  callback: commandMap,
    },
    "mapb": {
  name: "mapb",
  description: "Returns the previous 20 locations",
  callback: commandMapb,
    },
  }
}

func commandHelp(cfg *config) error {
  fmt.Println("Help: Use 'exit' to quit the program")
  return nil
}

func commandExit(cfg *config) error {
  fmt.Println("Quitting")
  os.Exit(0)
  return nil
}

func commandMap(cfg *config) error {
  var url string
  if len(cfg.next) > 0 {
    url = cfg.next
  } else {
    url = "https://pokeapi.co/api/v2/location" 
  }
  res, err := http.Get(url)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return nil
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
    fmt.Printf("body: %v\n", body)
  }

  response := Locations{}
  json.Unmarshal(body, &response)

  for v := range len(response.Results) {
    fmt.Printf("Location: %v\n", response.Results[v].Name)
  }

  fmt.Printf("response.Next: %v\n", response.Next)

  cfg.next = response.Next
  switch v := response.Previous.(type) {
  case string:
    cfg.previous = v
  }

  return nil
}

func commandMapb(cfg *config) error {
  var url string
  if len(cfg.previous) > 0 {
    url = cfg.previous
  } else {
    fmt.Println("No previous locations")
    return nil
  }
  res, err := http.Get(url)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return nil
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode > 299 {
    fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
    fmt.Printf("body: %v\n", body)
  }

  response := Locations{}
  json.Unmarshal(body, &response)

  for v := range len(response.Results) {
    fmt.Printf("Location: %v\n", response.Results[v].Name)
  }

  cfg.next = response.Next
  switch v := response.Previous.(type) {
  case string:  
    cfg.previous = v
  }

  return nil

}
