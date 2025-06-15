package main

import (
	"fmt"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func commandMap(s *replState, _ []string) (string, error) {
	if s.config.Next == nil {
		return "", fmt.Errorf("you are on the last page")
	}

	message, err := pokeapi.GetLocationAreas(s.config, s.config.Next, s.cache)
	if err != nil {
		return "", err
	}

	return message, nil
}
