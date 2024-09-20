package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/api/pokeapi"
	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokeclient"
)

type State struct {
	Previous *string
	Next     *string
}

type command struct {
	name        string
	description string
	callback    callbackFunc
}

type callbackFunc func(args []string) error

type pokedex map[string]pokeapi.Pokemon

var dexter = make(pokedex)

func main() {
	run()
}

func run() {
	client := pokeclient.NewClient(
		5*time.Minute,
		10*time.Second,
	)

	var state State

	commandMap := map[string]command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exitFunc,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    nil,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations in the Pokemon world",
			callback:    mapFunc(client, &state),
		},
		"mapb": {
			name:        "map back",
			description: "Displays the previous 20 locations in the Pokemon world",
			callback:    mapBFunc(client, &state),
		},
		"explore": {
			name:        "explore",
			description: "Lists all the Pokemon in a given area",
			callback:    exploreFunc(client),
		},
		"catch": {
			name:        "catch",
			description: "Catches a Pokemon and adds them to your Pokedex",
			callback:    catchFunc(client),
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a Pokemon from your Pokedex",
			callback:    inspectFunc(),
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists the names of all the Pokemon in your Pokedex",
			callback:    pokedexFunc(),
		},
	}

	summaries := summaryMap(commandMap)

	commandMap["help"] = command{
		name:        "help",
		description: "Displays a help message",
		callback:    helpFunc(summaries),
	}

	fmt.Printf("\nWelcome to the Pokedex!\n")
	fmt.Print("\npokedex > ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		command, args := parseArgs(input)

		cmd, ok := commandMap[command]
		if !ok {
			fmt.Println("ERROR: Unrecognised command.")

			fmt.Print("\npokedex > ")

			continue
		}

		if cmd.callback == nil {
			fmt.Println("ERROR: This command is defined but does not have a callback function.")

			fmt.Print("\npokedex > ")

			continue
		}

		if err := commandMap[command].callback(args); err != nil {
			fmt.Printf("ERROR: %v.\n", err)
		}

		fmt.Print("pokedex > ")
	}
}

func helpFunc(summaries map[string]string) callbackFunc {
	return func(_ []string) error {
		keys := []string{}

		for key := range maps.All(summaries) {
			keys = append(keys, key)
		}

		slices.Sort(keys)

		fmt.Printf("\nCommands:\n")

		for _, key := range slices.All(keys) {
			fmt.Printf("\n%s: %s", key, summaries[key])
		}

		fmt.Printf("\n\n")

		return nil
	}
}

func exitFunc(_ []string) error {
	os.Exit(0)

	return nil
}

func mapFunc(client *pokeclient.Client, state *State) callbackFunc {
	return func(_ []string) error {
		url := state.Next
		if url == nil {
			url = new(string)
			*url = pokeclient.LocationAreaPath
		}

		return printResourceList(client, *url, state)
	}
}

func mapBFunc(client *pokeclient.Client, state *State) callbackFunc {
	return func(_ []string) error {
		url := state.Previous
		if url == nil {
			return fmt.Errorf("no previous locations available")
		}

		return printResourceList(client, *url, state)
	}
}

func exploreFunc(client *pokeclient.Client) callbackFunc {
	return func(args []string) error {
		if args == nil {
			return errors.New("the location has not been specified")
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"unexpected number of locations: want 1; got %d",
				len(args),
			)
		}

		location := args[0]

		fmt.Println("Exploring", location)

		locationArea, err := client.GetLocationArea(location)
		if err != nil {
			return fmt.Errorf(
				"unable to get the location area: %w",
				err,
			)
		}

		fmt.Println("Found Pokemon:")

		for _, encounter := range slices.All(locationArea.PokemonEncounters) {
			fmt.Printf("- %s\n", encounter.Pokemon.Name)
		}

		return nil
	}
}

func catchFunc(client *pokeclient.Client) callbackFunc {
	return func(args []string) error {
		if args == nil {
			return errors.New("the name of the Pokemon has not been specified")
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"unexpected number of Pokemon names: want 1; got %d",
				len(args),
			)
		}

		pokemonName := args[0]

		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

		pokemon, err := client.GetPokemon(pokemonName)
		if err != nil {
			return fmt.Errorf(
				"unable to get the information on %s: %w",
				pokemonName,
				err,
			)
		}

		chance := 50

		if caught := catchPokemon(chance); caught {
			dexter[pokemonName] = pokemon
			fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", pokemonName)
		} else {
			fmt.Printf("%s escaped!\n", pokemonName)
		}

		return nil
	}
}

func inspectFunc() callbackFunc {
	return func(args []string) error {
		if args == nil {
			return errors.New("the name of the Pokemon has not been specified")
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"unexpected number of Pokemon names: want 1; got %d",
				len(args),
			)
		}

		pokemonName := args[0]

		pokemon, ok := dexter[pokemonName]
		if !ok {
			return fmt.Errorf("you have not caught %s", pokemonName)
		}

		info := fmt.Sprintf(
			"Name: %s\nHeight: %d\nWeight: %d\nStats:",
			pokemon.Name,
			pokemon.Height,
			pokemon.Weight,
		)

		for _, stat := range slices.All(pokemon.Stats) {
			info += fmt.Sprintf(
				"\n  - %s: %d",
				stat.Stat.Name,
				stat.BaseStat,
			)
		}

		info += "\nTypes:"

		for _, pType := range slices.All(pokemon.Types) {
			info += "\n  - " + pType.Type.Name
		}

		fmt.Println(info)

		return nil
	}
}

func pokedexFunc() callbackFunc {
	return func(_ []string) error {
		if len(dexter) == 0 {
			fmt.Println("You have no Pokemon in your Pokedex")

			return nil
		}

		fmt.Println("Your Pokedex:")

		for name := range maps.All(dexter) {
			fmt.Println("  -", name)
		}

		return nil
	}
}

func printResourceList(client *pokeclient.Client, url string, state *State) error {
	list, err := client.GetNamedAPIResourceList(url)
	if err != nil {
		return fmt.Errorf("unable to get the list of resources: %w", err)
	}

	state.Next = list.Next
	state.Previous = list.Previous

	for _, location := range slices.All(list.Results) {
		fmt.Println(location.Name)
	}

	return nil
}

func summaryMap(commandMap map[string]command) map[string]string {
	summaries := make(map[string]string)

	for key, value := range maps.All(commandMap) {
		summaries[key] = value.description
	}

	return summaries
}

func parseArgs(input string) (string, []string) {
	split := strings.Split(input, " ")

	if len(split) == 0 {
		return "", nil
	}

	if len(split) == 1 {
		return split[0], nil
	}

	return split[0], split[1:]
}

func catchPokemon(chance int) bool {
	if chance >= 100 {
		return true
	}

	if chance <= 0 {
		return false
	}

	maxInt := 100

	numGenerator := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	luckyNumberSet := make(map[int]struct{})

	for len(luckyNumberSet) < chance {
		num := numGenerator.IntN(maxInt)
		if _, ok := luckyNumberSet[num]; !ok {
			luckyNumberSet[num] = struct{}{}
		}
	}

	roller := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	got := roller.IntN(maxInt)

	_, ok := luckyNumberSet[got]

	return ok
}
