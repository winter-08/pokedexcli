package main

import (
	"bufio"
	"fmt"
	"math/rand"
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

var pokedex map[string]*pokedexapi.Pokemon = make(map[string]*pokedexapi.Pokemon)

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
    "catch": {
  name: "catch",
  description: "Returns the area for a given pokemon",
  callback: commandCatch,
    },
    "inspect": {
  name: "inspect",
  description: "Returns the area for a given pokemon",
  callback: commandInspect,
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
    fmt.Printf("There was an error getting locations:\n%v\n", err)
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
  if len(args) == 0 {
    fmt.Printf("Please give a location to explore")
    return nil
  }

  response, err := pokedexapi.GetLocationArea(args[0], cache)
  if err != nil {
    fmt.Printf("There was an error getting location area:\n%v\n", err)
    return nil
  }

  fmt.Printf("Pokemon:\n")
  for v := range response.PokemonEncounters {
    fmt.Printf("%v\n", response.PokemonEncounters[v].Pokemon.Name)
  }

  return nil
}

func commandCatch(cfg *config, cache *pokecache.Cache, args []string) error {
  if len(args) == 0 {
    fmt.Printf("Please give a pokemon to catch")
    return nil
  }

  response, err := pokedexapi.GetPokemon(args[0], cache)
  if err != nil {
    fmt.Printf("err: %v\n", err)
    return nil
  }

  fmt.Printf("Throwing a ball at %v...\n", response.Name)

  if 50 / rand.Intn(response.BaseExperience) > 1 {
    fmt.Printf("%v was caught!\n", response.Name)
    pokedex[response.Name] = response
  } else {
    fmt.Printf("%v escaped!\n", response.Name)
  }
  return nil  
}

func commandInspect(cfg *config, cache *pokecache.Cache, args []string) error {
  if len(args) == 0 {
    fmt.Printf("Please give a pokemon to inspect\n")
    return nil
  }

  pokemon, ok := pokedex[args[0]]

  if !ok {
    fmt.Print("you have not caught that pokemon\n")
    return nil
  }

  fmt.Printf("Name: %v\n", pokemon.Name)
  fmt.Printf("Height: %v\n", pokemon.Height)
  fmt.Printf("Weight: %v\n", pokemon.Weight)
  fmt.Print("Stats:\n")
  for v := range (pokemon.Stats) {
    fmt.Printf("  -%v: %d\n", pokemon.Stats[v].Stat.Name, pokemon.Stats[v].BaseStat)
  }
  fmt.Print("Types:\n")
  for v := range (pokemon.Types) {
    fmt.Printf("  - %v\n", pokemon.Types[v].Type.Name)
  }

  return nil
}
