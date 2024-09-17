package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
	"time"
)

const (
	baseURL              string = "https://pokeapi.co/api/v2"
	locationAreaEndpoint string = "/location-area"
)

type State struct {
	Previous *string
	Next     *string
}

var state State

type command struct {
	name        string
	description string
	callback    func() error
}

func main() {
	run()
}

func run() {
	fmt.Print("pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()

		cmdMap := getCommandMap()
		if _, ok := cmdMap[command]; !ok {
			fmt.Println("ERROR: Unrecognised command.")

			fmt.Print("\npokedex > ")

			continue
		}

		if err := cmdMap[command].callback(); err != nil {
			fmt.Printf("ERROR: %v.\n", err)
		}

		fmt.Print("pokedex > ")
	}
}

func getCommandMap() map[string]command {
	return map[string]command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Displays the previous 20 locations in the Pokemon world",
			callback:    commandMapB,
		},
	}
}

func commandHelp() error {
	cmdMap := getCommandMap()

	keys := []string{}

	for key := range maps.All(cmdMap) {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n")

	for _, key := range slices.All(keys) {
		fmt.Printf("\n%s: %s", key, cmdMap[key].description)
	}

	fmt.Println("\n")

	return nil
}

func commandExit() error {
	os.Exit(0)

	return nil
}

func commandMap() error {
	url := state.Next
	if url == nil {
		url = new(string)
		*url = baseURL + locationAreaEndpoint
	}

	return printMap(*url)
}

func commandMapB() error {
	url := state.Previous
	if url == nil {
		return fmt.Errorf("no previous locations available")
	}

	return printMap(*url)
}

func printMap(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating the HTTP request: %w", err)
	}

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error getting the response from the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf(
			"received a bad status from %s: (%d) %s",
			url,
			resp.StatusCode,
			resp.Status,
		)
	}

	var result NamedAPIResourceList

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("unable to decode the JSON response: %w", err)
	}

	state.Next = result.Next
	state.Previous = result.Previous

	for _, location := range slices.All(result.Results) {
		fmt.Println(location.Name)
	}

	return nil
}
