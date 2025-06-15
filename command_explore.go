package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
)

func commandExplore(s *replState, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("you must provide a location area name")
	}
	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName + "/"

	data, exists := s.cache.Get(url)
	if !exists {
		response, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("error fetching data: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode > 299 {
			return "", fmt.Errorf("bad status code: %d - location not found", response.StatusCode)
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error reading data: %w", err)
		}
		s.cache.Add(url, data)
	}

	var locationAreaInfo pokeapi.LocationAreaInfo
	err := json.Unmarshal(data, &locationAreaInfo)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling data: %w", err)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Exploring %s...\n", locationAreaInfo.Name))
	builder.WriteString("Found Pokémon:\n")
	for _, encounter := range locationAreaInfo.PokemonEncounters {
		builder.WriteString(fmt.Sprintf(" - %s\n", encounter.Pokemon.Name))
	}

	if exists {
		builder.WriteString("--------------------------------\n")
		builder.WriteString("Cache hit!\n")
	}
	builder.WriteString("--------------------------------\n")

	return builder.String(), nil
}
