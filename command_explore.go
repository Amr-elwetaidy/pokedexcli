package main

import (
	"fmt"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func commandExplore(s *replState, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("you must provide a location area name")
	}
	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName + "/"

	message, err := pokeapi.GetLocationArea(url, s.cache)
	if err != nil {
		return "", err
	}

	return message, nil
}
