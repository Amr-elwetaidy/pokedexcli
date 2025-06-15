package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
	"github.com/Amr-elwetaidy/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	re := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	cleaned := re.ReplaceAllString(text, "")
	lowercase := strings.ToLower(cleaned)
	words := strings.Fields(lowercase)
	return words
}

func commandExit(_ *replState) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *replState, commands map[string]cliCommand) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s - %s\n", command.name, command.description)
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func main() {
	state := &replState{
		config: &pokeapi.Config{
			Next: strPtr("https://pokeapi.co/api/v2/location-area/"),
		},
		cache: pokecache.NewCache(10 * time.Second),
	}

	commands := make(map[string]cliCommand)

	commands["map"] = cliCommand{
		name:        "map",
		description: "Show the next 20 location areas",
		callback:    func(s *replState) error { return pokeapi.CommandMap(s.config, s.cache) },
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Show the previous 20 location areas",
		callback:    func(s *replState) error { return pokeapi.CommandMapBack(s.config, s.cache) },
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Show all commands",
		callback: func(s *replState) error {
			return commandHelp(s, commands)
		},
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

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

		err := commandCallback.callback(state)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
