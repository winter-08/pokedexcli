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
