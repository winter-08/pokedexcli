package main

import (
	"bufio"
	"fmt"
	"pokedexcli/internal/pokedexapi"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

type config struct {
  next string 
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(*config, *pokecache.Cache) error
}


func main() {

  config := config {
    next: "",
    previous: "",
  }

  cache, err := pokecache.NewCache(5 * time.Minute)

  if err != nil {
    fmt.Printf("Error starting cache")
    return 
  }

  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := scanner.Text()


    commands := cliCommands()

    if cmd, ok := commands[input]; ok {
      cmd.callback(&config, cache)
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

func commandHelp(cfg *config, cache *pokecache.Cache) error {
  fmt.Println("Help: Use 'exit' to quit the program")
  return nil
}

func commandExit(cfg *config, cache *pokecache.Cache) error {
  fmt.Println("Quitting")
  os.Exit(0)
  return nil
}

func commandMap(cfg *config, cache *pokecache.Cache) error {
  var url string
  if len(cfg.next) > 0 {
    url = cfg.next
  } else {
    url = "https://pokeapi.co/api/v2/location" 
  }
  response, err := pokedexapi.GetLocations(url, cache)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return nil
  }

  for v := range response.Results {
    fmt.Printf("Location: %v\n", response.Results[v].Name)
  }


  cfg.next = response.Next
  switch v := response.Previous.(type) {
  case string:
    cfg.previous = v
  }

  return nil
}

func commandMapb(cfg *config, cache *pokecache.Cache) error {
  var url string
  if len(cfg.previous) > 0 {
    url = cfg.previous
  } else {
    fmt.Println("No previous locations")
    return nil
  }
  response, err := pokedexapi.GetLocations(url, cache)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return nil
  }

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
