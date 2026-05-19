package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pokecache "github.com/Shamiul-isl/pokedexcli/internal"
)

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      &confStruct{},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      &confStruct{},
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    commandMap,
			config:      &confStruct{},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
			config:      &confStruct{},
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of the pokemon found in this location area",
			callback:    commandExplore,
			config:      &confStruct{},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
			config:      &confStruct{},
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect caught pokemon",
			callback:    commandInspect,
			config:      &confStruct{},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all caught pokemon in pokedex",
			callback:    commandPokedex,
			config:      &confStruct{},
		},
	}
	Pokedex = make(map[string]Pokemon)
	cache = pokecache.NewCache(5 * time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := strings.ToLower(strings.TrimSpace(scanner.Text()))
			fields := strings.Fields(text)
			if value, ok := commands[fields[0]]; !ok {
				fmt.Println("Unknown command")
				continue
			} else {
				// param := ""
				if len(fields) < 2 {
					value.callback("")
				} else {
					value.callback(fields[1])
				}
			}
		}
	}
}
