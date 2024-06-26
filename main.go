package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"pokedexcli/internal/pokedexapi"
	"strings"
	"time"
)

type config struct {
  next string 
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(*config, *pokecache.Cache, []string) error
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
  
    args := strings.Split(input, " ")

    commands := cliCommands()

    if cmd, ok := commands[args[0]]; ok {
      cmd.callback(&config, cache, args[1:])
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
    "explore": {
  name: "explore",
  description: "Returns the area for a given location",
  callback: commandExplore,
    },
  }
}

func commandHelp(cfg *config, cache *pokecache.Cache, args []string) error {
  fmt.Println("Help: Use 'exit' to quit the program")
  return nil
}

func commandExit(cfg *config, cache *pokecache.Cache, args []string) error {
  fmt.Println("Quitting")
  os.Exit(0)
  return nil
}

func commandMap(cfg *config, cache *pokecache.Cache, args []string) error {
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

  fmt.Printf("Locations:\n")
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

func commandMapb(cfg *config, cache *pokecache.Cache, args []string) error {
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

  fmt.Printf("Locations:\n")
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

func commandExplore(cfg *config, cache *pokecache.Cache, args []string) error {
  fmt.Printf("args: %v\n", args)
  if len(args) == 0 {
    fmt.Printf("Please give a location to explore")
    return nil
  }

  response, err := pokedexapi.GetLocationArea(args[0], cache)
  if err != nil {
    fmt.Printf("There was an error getting location area: %v\n", err)
  }

  fmt.Printf("Pokemon:\n")
  for v := range response.PokemonEncounters {
    fmt.Printf("%v\n", response.PokemonEncounters[v].Pokemon.Name)
  }

  return nil
}
