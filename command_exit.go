package main

import (
	"fmt"
)

var ErrExit = fmt.Errorf("exit")

func commandExit(_ *replState, _ []string) (string, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return "", ErrExit
}
