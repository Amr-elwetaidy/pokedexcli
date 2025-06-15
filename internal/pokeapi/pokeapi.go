package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokecache"
)

func GetLocationAreas(config *Config, url *string, cache *pokecache.Cache) error {
	if url == nil {
		return fmt.Errorf("cannot fetch location areas, URL is nil")
	}

	var locationAreas LocationAreas
	data, exists := cache.Get(*url)
	if !exists {
		response, err := http.Get(*url)
		if err != nil {
			return fmt.Errorf("error fetching data: %w", err)
		}
		defer response.Body.Close()

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading data: %w", err)
		}
		cache.Add(*url, data)
	}

	err := json.Unmarshal(data, &locationAreas)
	if err != nil {
		return fmt.Errorf("error unmarshalling data: %w", err)
	}

	config.Next = &locationAreas.Next
	config.Previous = locationAreas.Previous

	printLocationAreas(locationAreas)
	if exists {
		fmt.Println("Cache hit!")
	}
	fmt.Println("--------------------------------")

	return nil
}

func CommandMap(config *Config, cache *pokecache.Cache) error {
	if config.Next == nil {
		fmt.Println("You're on the last page")
		return nil
	}
	err := GetLocationAreas(config, config.Next, cache)
	if err != nil {
		return err
	}

	return nil
}

func CommandMapBack(config *Config, cache *pokecache.Cache) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	err := GetLocationAreas(config, config.Previous, cache)
	if err != nil {
		return err
	}

	return nil
}

func printLocationAreas(locationAreas LocationAreas) {
	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.GetName())
	}
}
