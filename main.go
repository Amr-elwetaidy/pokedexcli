package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
}

func commandExit(_ *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand, _ *pokeapi.Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s - %s\n", command.name, command.description)
	}
	return nil
}

func main() {
	var commands map[string]cliCommand
	commands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Show the next 20 location areas",
			callback:    pokeapi.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 location areas",
			callback:    pokeapi.CommandMapBack,
		},
		"help": {
			name:        "help",
			description: "Show all commands",
			callback:    func(config *pokeapi.Config) error { return commandHelp(commands, config) },
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	next := "https://pokeapi.co/api/v2/location-area/"
	config := &pokeapi.Config{Next: &next}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		command := scanner.Text()

		words := cleanInput(command)
		if len(words) == 0 {
			continue
		}

		commandCallback, exists := commands[words[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := commandCallback.callback(config)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
