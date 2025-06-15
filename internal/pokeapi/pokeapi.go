package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokecache"
)

func GetLocationAreas(config *Config, url *string, cache *pokecache.Cache) (string, error) {
	if url == nil {
		return "", fmt.Errorf("cannot fetch location areas, URL is nil")
	}

	var locationAreas LocationAreas
	data, exists := cache.Get(*url)
	if !exists {
		response, err := http.Get(*url)
		if err != nil {
			return "", fmt.Errorf("error fetching data: %w", err)
		}
		defer response.Body.Close()

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error reading data: %w", err)
		}
		cache.Add(*url, data)
	}

	err := json.Unmarshal(data, &locationAreas)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling data: %w", err)
	}

	config.Next = &locationAreas.Next
	config.Previous = locationAreas.Previous

	var builder strings.Builder
	printLocationAreas(locationAreas, &builder)
	if exists {
		builder.WriteString("--------------------------------\n")
		builder.WriteString("Cache hit!\n")
	}
	builder.WriteString("--------------------------------\n")

	return builder.String(), nil
}

func printLocationAreas(locationAreas LocationAreas, builder *strings.Builder) {
	for _, locationArea := range locationAreas.Results {
		builder.WriteString(locationArea.GetName() + "\n")
	}
}

func GetLocationArea(url string, cache *pokecache.Cache) (string, error) {
	data, exists := cache.Get(url)
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
		cache.Add(url, data)
	}

	var locationAreaInfo LocationAreaInfo
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

func GetPokemon(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"
	var pokemon Pokemon
	data, exists := cache.Get(url)
	if !exists {
		resp, err := http.Get(url)
		if err != nil {
			return Pokemon{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return Pokemon{}, err
		}
		cache.Add(url, data)
	}

	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}
