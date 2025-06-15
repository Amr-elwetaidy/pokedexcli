package main

import (
	"fmt"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func commandMapb(s *replState, _ []string) (string, error) {
	if s.config.Previous == nil {
		return "", fmt.Errorf("you are on the first page")
	}

	message, err := pokeapi.GetLocationAreas(s.config, s.config.Previous, s.cache)
	if err != nil {
		return "", err
	}

	return message, nil
}
