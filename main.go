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

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokedex/internal/poketrainer"
)

type callbackFunc func(args []string) error

type command struct {
	name        string
	description string
	callback    callbackFunc
}

func main() {
	run()
}

func run() {
	client := pokeclient.NewClient(
		5*time.Minute,
		10*time.Second,
	)

	trainer := poketrainer.NewTrainer()

	commandMap := map[string]command{
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon and add it to your Pokedex",
			callback:    catchFunc(client, trainer),
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exitFunc,
		},
		"explore": {
			name:        "explore",
			description: "List all the Pokemon in a given area",
			callback:    exploreFunc(client, trainer),
		},
		"help": {
			name:        "help",
			description: "Display the help message",
			callback:    nil,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon from your Pokedex",
			callback:    inspectFunc(trainer),
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations in the Pokemon world",
			callback:    mapFunc(client, trainer),
		},
		"mapb": {
			name:        "map back",
			description: "Display the previous 20 locations in the Pokemon world",
			callback:    mapBFunc(client, trainer),
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the names of all the Pokemon in your Pokedex",
			callback:    pokedexFunc(trainer),
		},
		"visit": {
			name:        "visit",
			description: "Visit a location area",
			callback:    visitFunc(client, trainer),
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

func mapFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) callbackFunc {
	return func(_ []string) error {
		url := trainer.NextLocationArea()
		if url == nil {
			url = new(string)
			*url = pokeclient.LocationAreaPath
		}

		return printResourceList(client, *url, trainer.UpdateLocationAreas)
	}
}

func mapBFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) callbackFunc {
	return func(_ []string) error {
		url := trainer.PreviousLocationArea()
		if url == nil {
			return fmt.Errorf("no previous locations available")
		}

		return printResourceList(client, *url, trainer.UpdateLocationAreas)
	}
}

func exploreFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) callbackFunc {
	return func(_ []string) error {
		locationAreaName := trainer.CurrentLocationAreaName()

		fmt.Printf("Exploring %s...\n", locationAreaName)

		locationArea, err := client.GetLocationArea(locationAreaName)
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

func visitFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) callbackFunc {
	return func(args []string) error {
		if args == nil {
			return errors.New("the location area has not been specified")
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"unexpected number of location areas: want 1; got %d",
				len(args),
			)
		}

		locationAreaName := args[0]

		locationArea, err := client.GetLocationArea(locationAreaName)
		if err != nil {
			return fmt.Errorf(
				"unable to get the location area: %w",
				err,
			)
		}

		trainer.UpdateCurrentLocationAreaName(locationArea.Name)

		fmt.Println("You are now visiting", locationArea.Name)

		return nil
	}
}

func catchFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) callbackFunc {
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

		pokemonDetails, err := client.GetPokemon(pokemonName)
		if err != nil {
			return fmt.Errorf(
				"unable to get the information on %s: %w",
				pokemonName,
				err,
			)
		}

		encountersPath := pokemonDetails.LocationAreaEncounters

		encounterAreas, err := client.GetPokemonLocationAreas(encountersPath)
		if err != nil {
			return fmt.Errorf(
				"unable to get the Pokemon's possible encounter areas: %w",
				err,
			)
		}

		validLocationArea := false
		currentLocation := trainer.CurrentLocationAreaName()

		for _, area := range slices.All(encounterAreas) {
			if currentLocation == area.LocationArea.Name {
				validLocationArea = true

				break
			}
		}

		if !validLocationArea {
			return fmt.Errorf(
				"%s cannot be found in %s",
				pokemonName,
				currentLocation,
			)
		}

		chance := 50

		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

		if caught := catchPokemon(chance); caught {
			trainer.AddPokemonToPokedex(pokemonName, pokemonDetails)
			fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", pokemonName)
		} else {
			fmt.Printf("%s escaped!\n", pokemonName)
		}

		return nil
	}
}

func inspectFunc(trainer *poketrainer.Trainer) callbackFunc {
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

		pokemon, ok := trainer.GetPokemonFromPokedex(pokemonName)
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

func pokedexFunc(trainer *poketrainer.Trainer) callbackFunc {
	return func(_ []string) error {
		trainer.ListAllPokemonFromPokedex()

		return nil
	}
}

func printResourceList(
	client *pokeclient.Client,
	url string,
	updateStateFunc func(previous *string, next *string),
) error {
	list, err := client.GetNamedAPIResourceList(url)
	if err != nil {
		return fmt.Errorf("unable to get the list of resources: %w", err)
	}

	if updateStateFunc != nil {
		updateStateFunc(list.Previous, list.Next)
	}

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
