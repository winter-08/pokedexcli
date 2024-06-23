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
  fmt.Printf("Help: Use 'exit' to quit the program\n")
  return nil
}

func commandExit() error {
  fmt.Printf("Quitting")
  os.Exit(0)
  return nil
}
