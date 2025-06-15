package main

import (
	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
	"github.com/Amr-elwetaidy/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*replState, []string) (string, error)
}

type replState struct {
	config *pokeapi.Config
	cache  *pokecache.Cache
}
