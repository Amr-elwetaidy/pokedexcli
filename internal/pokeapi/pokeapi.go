package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CommandMap(config *Config) error {
	if config.Next == nil {
		fmt.Println("You're on the last page")
		return nil
	}

	response, err := http.Get(*config.Next)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	locationAreas := LocationAreas{}
	err = json.NewDecoder(response.Body).Decode(&locationAreas)
	if err != nil {
		return err
	}

	*config.Next = locationAreas.Next
	config.Previous = locationAreas.Previous

	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.GetName())
	}

	return nil
}

func CommandMapBack(config *Config) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	response, err := http.Get(*config.Previous)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	locationAreas := LocationAreas{}
	err = json.NewDecoder(response.Body).Decode(&locationAreas)
	if err != nil {
		return err
	}

	*config.Next = locationAreas.Next
	config.Previous = locationAreas.Previous

	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.GetName())
	}

	return nil
}
