package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
  name string
  description string
  callback func() error
}

func main() {

  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := scanner.Text()

    commands := cliCommands()

    if cmd, ok := commands[input]; ok {
      cmd.callback()
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

func commandHelp() error {
  fmt.Println("Help: Use 'exit' to quit the program")
  return nil
}

func commandExit() error {
  fmt.Println("Quitting")
  os.Exit(0)
  return nil
}

func commandMap() error {
    return nil
}

func commandMapb() error {
    return nil
}
