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
