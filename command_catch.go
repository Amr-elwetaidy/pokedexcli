package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func commandCatch(s *replState, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, err := pokeapi.GetPokemon(pokemonName, s.cache)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Throwing a Pokeball at %s...\n", pokemon.Name))

	// The higher the base experience, the harder to catch.
	const catchThreshold = 60
	catchValue := rand.Intn(pokemon.BaseExperience)

	if catchValue > catchThreshold {
		builder.WriteString(fmt.Sprintf("%s escaped!", pokemon.Name))
		return builder.String(), nil
	}

	builder.WriteString(fmt.Sprintf("%s was caught!", pokemon.Name))
	builder.WriteString("\n(You may now inspect it with the 'inspect' command)")
	s.pokedex[pokemonName] = pokemon

	return builder.String(), nil
}
