package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	pokecache "github.com/Shamiul-isl/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
	config      *confStruct
}

type confStruct struct {
	Next     string
	Previous string
}

type areaResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type areaBody struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []areaResult `json:"results"`
}

type encounterBody struct {
	Encounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon areaResult `json:"pokemon"`
}

type pokemonBody struct {
	Exp    int        `json:"base_experience"`
	Height int        `json:"height"`
	Stats  []statBody `json:"stats"`
	Types  []typeBody `json:"types"`
	Weight int        `json:"weight"`
}

type statBody struct {
	Base_Stat int        `json:"base_stat"`
	Stat      areaResult `json:"stat"`
}

type typeBody struct {
	Type areaResult `json:"type"`
}

type Pokemon struct {
	Name     string
	Catchpct float64
	RespBody pokemonBody
}

var commands map[string]cliCommand
var cache pokecache.Cache
var Pokedex map[string]Pokemon

func commandExit(param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

func commandHelp(param string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key := range commands {
		fmt.Printf("%s: %s\n", commands[key].name, commands[key].description)
	}
	return nil
}

func commandMap(param string) error {
	urlToUse := "https://pokeapi.co/api/v2/location-area/"
	if commands["map"].config.Next != "" {
		urlToUse = commands["map"].config.Next
	}

	var data []byte
	var found bool

	if data, found = cache.Get(urlToUse); !found {
		res, err := http.Get(urlToUse)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		cache.Add(urlToUse, data)
	}

	var areas areaBody
	if err := json.Unmarshal(data, &areas); err != nil {
		return err
	}
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	commands["map"].config.Previous = areas.Previous
	commands["map"].config.Next = areas.Next
	return nil
}

func commandMapB(param string) error {
	if commands["map"].config.Previous == "" {
		fmt.Println("you're on the first page")
		return fmt.Errorf("you're on the first page")
	}
	urlToUse := commands["map"].config.Previous
	var data []byte
	var found bool

	if data, found = cache.Get(urlToUse); !found {
		res, err := http.Get(urlToUse)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		cache.Add(urlToUse, data)
	}

	var areas areaBody
	if err := json.Unmarshal(data, &areas); err != nil {
		return err
	}
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	commands["map"].config.Previous = areas.Previous
	commands["map"].config.Next = areas.Next
	return nil
}

func commandExplore(location_area string) error {
	urlToUse := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location_area)

	var data []byte
	var found bool

	if data, found = cache.Get(urlToUse); !found {
		res, err := http.Get(urlToUse)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		cache.Add(urlToUse, data)
	}

	var enc encounterBody
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", location_area)
	fmt.Println("Found Pokemon:")
	for _, poke := range enc.Encounters {
		fmt.Printf("- %s\n", poke.Pokemon.Name)
	}
	return nil
}

func commandCatch(pokename string) error {
	urlToUse := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokename)

	var data []byte
	var found bool

	if data, found = cache.Get(urlToUse); !found {
		res, err := http.Get(urlToUse)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		cache.Add(urlToUse, data)
	}

	var enc pokemonBody
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokename)

	choose := rand.Float64()
	if choose > (float64(enc.Exp) / 306.0) {
		fmt.Printf("%s was caught!\n", pokename)
		Pokedex[pokename] = Pokemon{
			Name:     pokename,
			Catchpct: float64(enc.Exp) / 306.0,
			RespBody: enc,
		}
	} else {
		fmt.Printf("%s escaped!\n", pokename)
	}

	// fmt.Println("Found Pokemon:")
	return nil
}

func commandInspect(pokename string) error {
	if _, found := Pokedex[pokename]; !found {
		fmt.Println("you have not caught that pokemon")
		return errors.New("you have not caught that pokemon")
	}

	val := Pokedex[pokename].RespBody

	fmt.Printf("Name: %s\n", pokename)
	fmt.Printf("Height: %d\n", val.Height)
	fmt.Printf("Weight: %d\n", val.Weight)
	fmt.Println("Stats:")
	for _, s := range val.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.Base_Stat)
	}
	fmt.Println("Types:")
	for _, t := range val.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(pokename string) error {
	fmt.Println("Your Pokedex:")
	for key := range Pokedex {
		fmt.Printf("  - %s\n", key)
	}
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
