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
  stateMap
  statePokedex
  stateCatch
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
  commands map[string]cliCommand
}

type config struct {
  next string 
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(model, tea.Msg) (model,tea.Cmd, error)
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
    commands: cliCommands(),
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
  case stateMap:
    return m.updateMap()
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

        if cmd, ok := m.commands[m.choices[m.cursor]]; ok {
          cmd.callback(m, msg)
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
      m, cmd, _ := m.executeExplore()
      return m, cmd 
    case tea.KeyCtrlC, tea.KeyEsc:
      m.state = stateMainMenu
      return m, nil 
    }
    
  }

  m.textInput, cmd = m.textInput.Update(msg)

  return m, cmd 
}

func (m model) executeExplore() (tea.Model, tea.Cmd, error) {
  location := strings.TrimSpace(m.textInput.Value())

  response, err := pokedexapi.GetLocationArea(location, m.cache)
  if err != nil {
    m.err = fmt.Errorf("error getting location area %v", err)
    return m, nil, nil
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
  return m, nil, nil

}

func (m model) updateMap() (tea.Model, tea.Cmd) {
  var url string
  if len(m.config.next) > 0 {
    url = m.config.next
  } else {
    url = "https://pokeapi.co/api/v2/location" 
  }
  response, err := pokedexapi.GetLocations(url, m.cache)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return m, nil
  }

  rows := []table.Row{}
  for _, location := range response.Results {
    rows = append(rows, table.Row{location.Name, location.Region.Name})
  }

  columns := []table.Column{
    { Title: "Location", Width: 20 },
    { Title: "Region", Width: 20 },
  }

  t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(10),
  )

  m.table = t

  m.config.next = response.Next
  switch v := response.Previous.(type) {
  case string:
    m.config.previous = v
  }

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

func (m model) viewMap() string {
  return ""
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

func commandHelp(m model, msg tea.Msg) (model, tea.Cmd, error) {
  fmt.Println("Help: Use 'exit' to quit the program")
  return m, nil, nil
}

func commandExit(m model, msg tea.Msg) (model, tea.Cmd, error) {
  fmt.Println("Quitting")
  os.Exit(0)
  return m, nil, nil
}

func commandMap(m model, msg tea.Msg) (model, tea.Cmd, error) {
  var url string
  if len(m.config.next) > 0 {
    url = m.config.next
  } else {
    url = "https://pokeapi.co/api/v2/location" 
  }
  response, err := pokedexapi.GetLocations(url, m.cache)
  if err != nil {
    fmt.Printf("There was an error getting locations: %v\n", err)
    return m, nil, nil
  }

  fmt.Printf("Locations:\n")
  for v := range response.Results {
    fmt.Printf("Location: %v\n", response.Results[v].Name)
  }


  m.config.next = response.Next
  switch v := response.Previous.(type) {
  case string:
    m.config.previous = v
  }

  return m, nil, nil
}

func commandMapb(m model, msg tea.Msg) (model, tea.Cmd, error) {
  var url string
  if len(m.config.previous) > 0 {
    url = m.config.previous
  } else {
    fmt.Println("No previous locations")
    return m, nil, nil
  }
  response, err := pokedexapi.GetLocations(url, m.cache)
  if err != nil {
    fmt.Printf("There was an error getting locations:\n%v\n", err)
    return m, nil, nil
  }

  fmt.Printf("Locations:\n")
  for v := range len(response.Results) {
    fmt.Printf("Location: %v\n", response.Results[v].Name)
  }

  m.config.next = response.Next
  switch v := response.Previous.(type) {
  case string:  
    m.config.previous = v
  }

  return m, nil, nil

}

write a hello world function

func commandExplore(m model, msg tea.Msg) (model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyEnter:
      m.state = stateMainMenu
      m, cmd := m.executeExplore()
      return m, cmd
    case tea.KeyCtrlC, tea.KeyEsc:
      m.state = stateMainMenu
      return m, nil
    }
  }

  m.textInput, cmd = m.textInput.Update(msg)

  return m, cmd
}

func commandCatch(m model, msg tea.Msg) (model, tea.Cmd, error) {

  msg2 := msg.(tea.KeyMsg)

  msg_string := msg2.String()

  if len(msg_string) == 0 {
    fmt.Printf("Please give a pokemon to catch")
    return m, nil, nil
  }

  response, err := pokedexapi.GetPokemon(msg_string, m.cache)
  if err != nil {
    fmt.Printf("err: %v\n", err)
    return m, nil, nil
  }

  fmt.Printf("Throwing a ball at %v...\n", response.Name)

  if 50 / rand.Intn(response.BaseExperience) > 1 {
    fmt.Printf("%v was caught!\n", response.Name)
    pokedex[response.Name] = *response
  } else {
    fmt.Printf("%v escaped!\n", response.Name)
  }
  return m, nil  
}

func commandInspect(m model, msg tea.Msg) (model, tea.Cmd, error) {
  msg2 := msg.(tea.KeyMsg)

  msg_string := msg2.String()

  if len(msg_string) == 0 {
    fmt.Printf("Please give a pokemon to inspect\n")
    return m, nil, nil
  }

  pokemon, ok := pokedex[msg_string]

  if !ok {
    fmt.Print("you have not caught that pokemon\n")
    return m, nil, nil
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

  return m, nil, nil
}

func commandPokedex(m model, msg tea.Msg) (model, tea.Cmd, error) {
  if len(pokedex) == 0 {
    fmt.Print("You have not caught any pokemon yet\n")
    return m, nil, nil
  }

  fmt.Print("Your Pokedex:\n")

  for v := range (pokedex) {
    fmt.Printf("  - %v\n", pokedex[v].Name)
  }

  return m, nil, nil

}
