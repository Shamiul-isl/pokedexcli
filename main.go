package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := strings.ToLower(strings.TrimSpace(scanner.Text()))
			words := strings.Fields(text)
			fmt.Printf("Your command was: %s\n", words[0])
		}
	}
}
