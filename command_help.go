package main

import (
	"fmt"
	"strings"
)

func commandHelp(_ *replState, commands map[string]cliCommand, _ []string) (string, error) {
	var builder strings.Builder
	builder.WriteString("\nWelcome to the Pokedex!\n")
	builder.WriteString("Usage:\n")
	builder.WriteString("--------------------------------\n")
	for _, command := range commands {
		builder.WriteString(fmt.Sprintf("%s: %s\n", command.name, command.description))
	}
	builder.WriteString("--------------------------------\n")
	return builder.String(), nil
}
