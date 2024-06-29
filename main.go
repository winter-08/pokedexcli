package main

import (
	//"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/internal/pokecache"
	"pokedexcli/internal/pokedexapi"
	"strings"

	//"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type appState int

const (
  stateMainMenu appState = iota
  stateExplore
)

type model struct {
  choices []string
  cursor int
  selected map[int]struct{}
  cache *pokecache.Cache; err error
  config config
  state appState
  textInput textinput.Model
  table table.Model
}

type config struct {
  next string 
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(*config, *pokecache.Cache, []string) error
}

var pokedex map[string]pokedexapi.Pokemon = make(map[string]pokedexapi.Pokemon)

func main() {


  p := tea.NewProgram(initialModel(), tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    fmt.Printf("oh no! an error: %v", err)
    os.Exit(1)
  }


}

func initialModel() model {
  cache, err := pokecache.NewCache(5 * time.Minute)

  if err != nil {
    fmt.Printf("Error starting cache")
  }

  ti := textinput.New()

  ti.Placeholder = "Enter location to explore"
  ti.Focus()

  return model{
    choices: []string{ "help", "explore", "map", "mapb", "exit", "catch", "inspect", "pokedex" },
    selected: make(map[int]struct{}),
    cache: cache,
    config: config {
      next: "",
      previous: "",
    },
    state: stateMainMenu,
    textInput: ti,
  }
}

func (m model) Init() tea.Cmd {
  return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch m.state {
  case stateMainMenu:
    return m.updateMainMenu(msg)
  case stateExplore:
    return m.updateExplore(msg) 
  }
  return m, nil

}

func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg:= msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {

    case "ctrl+c", "q":
      return m, tea.Quit

    case "up", "k":
      if m.cursor > 0 {
        m.cursor--
      }
    
    case "down", "j":
      if m.cursor < len(m.choices) - 1 {
        m.cursor++
      }

    case "enter", " ":
      _, ok := m.selected[m.cursor]
      if ok {
        delete(m.selected, m.cursor)
      } else {
        m.selected[m.cursor] = struct{}{}
        if m.choices[m.cursor] == "explore" {
          m.state = stateExplore
        }
        commands := cliCommands()

        var args []string

        if cmd, ok := commands[m.choices[m.cursor]]; ok {
          cmd.callback(&m.config, m.cache, args)
        }
      }
    }
  }
  return m, nil
}

func (m model) updateExplore(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyEnter:
      m.state = stateMainMenu
      return m.executeExplore()
    case tea.KeyCtrlC, tea.KeyEsc:
      m.state = stateMainMenu
      return m, nil
    }
    
  }

  m.textInput, cmd = m.textInput.Update(msg)

  return m, cmd 
}

func (m model) executeExplore() (tea.Model, tea.Cmd) {
  location := strings.TrimSpace(m.textInput.Value())

  response, err := pokedexapi.GetLocationArea(location, m.cache)
  if err != nil {
    m.err = fmt.Errorf("error getting location area %v", err)
    return m, nil
  }
  
  rows := []table.Row{}
  for _, encounter := range response.PokemonEncounters {
    rows = append(rows, table.Row{encounter.Pokemon.Name, encounter.Pokemon.URL})
  }

  columns := []table.Column{
    {Title: "Pokemon", Width: 20},
    {Title: "url", Width: 20},
  }
  t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(10),
  )

  m.table = t
  return m, nil

}

func (m model) View() string {
  clearScreen()
  switch m.state {
    case stateMainMenu:
      return m.viewMainMenu()
    case stateExplore:
      return m.viewExplore()
  }
  return ""
}

func (m model) viewMainMenu() string {
  s := "Pokedex"

  for i, choice := range m.choices {
    cursor := " "
    if m.cursor == i {
      cursor = ">"
    }

    checked := " "
    if _, ok := m.selected[i]; ok {
      checked = "x"
    }

    s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
  }

  s += "\nPress q to quit.\n"

  return s

}

func (m model) viewExplore() string {
  return fmt.Sprintf(
    "Enter location to explore: \n\n%s\n\n(Press Enter to confirm, Esc to cancel)",
    m.textInput.View(),
  )
}

func clearScreen() {
  fmt.Print("\033[23")
  fmt.Print("\033[H")
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
  description: "Attempts to catch a given pokemon and add to your pokedex",
  callback: commandCatch,
    },
    "inspect": {
  name: "inspect",
  description: "Inspects pokemon stats for a pokemon from your pokedex",
  callback: commandInspect,
    },
    "pokedex": {
  name: "pokedex",
  description: "Returns all pokemon in your pokedex",
  callback: commandPokedex,
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
    pokedex[response.Name] = *response
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

func commandPokedex(cfg *config, cache *pokecache.Cache, args []string) error {
  if len(pokedex) == 0 {
    fmt.Print("You have not caught any pokemon yet\n")
    return nil
  }

  fmt.Print("Your Pokedex:\n")

  for v := range (pokedex) {
    fmt.Printf("  - %v\n", pokedex[v].Name)
  }

  return nil

}
