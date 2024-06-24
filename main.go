package main

import (
	"bufio"
	"fmt"
	"os"
)

type config struct {
  next string
  previous string
}

type cliCommand struct {
  name string
  description string
  callback func(config) error
}

func main() {

  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := scanner.Text()

    config := config {
      next: "",
      previous: "",
    }

    commands := cliCommands(config)

    if cmd, ok := commands[input]; ok {
      cmd.callback(config)
    }

  }

}

func cliCommands(config) map[string]cliCommand {
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

func commandHelp(config) error {
  fmt.Println("Help: Use 'exit' to quit the program")
  return nil
}

func commandExit(config) error {
  fmt.Println("Quitting")
  os.Exit(0)
  return nil
}

func commandMap(config) error {
    return nil
}

func commandMapb(config) error {
    return nil
}
