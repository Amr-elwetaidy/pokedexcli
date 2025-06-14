package main

import "github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}
