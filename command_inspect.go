package main

import (
	"fmt"
)

func commandInspect(s *replState, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("you must provide a pokemon name")
	}

	pokemonName := args[0]
	pokemon, ok := s.pokedex[pokemonName]
	if !ok {
		return "you have not caught that pokemon", nil
	}

	var output string
	output += fmt.Sprintf("Name: %s\n", pokemon.Name)
	output += fmt.Sprintf("Height: %d\n", pokemon.Height)
	output += fmt.Sprintf("Weight: %d\n", pokemon.Weight)
	output += "Stats:\n"
	for _, stat := range pokemon.Stats {
		output += fmt.Sprintf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	output += "Types:\n"
	for _, t := range pokemon.Types {
		output += fmt.Sprintf("  - %s\n", t.Type.Name)
	}

	return output, nil
}
