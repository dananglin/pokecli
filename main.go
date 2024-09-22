package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"

	"codeflow.dananglin.me.uk/apollo/pokecli/internal/commands"
	"codeflow.dananglin.me.uk/apollo/pokecli/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokecli/internal/poketrainer"
)

type command struct {
	description string
	callback    commands.CommandFunc
}

func main() {
	run()
}

func run() {
	var (
		cacheCleanupInterval = 30 * time.Minute
		httpTimeout          = 10 * time.Second
		client               = pokeclient.NewClient(cacheCleanupInterval, httpTimeout)
		trainer              = poketrainer.NewTrainer()
	)

	commandMap := map[string]command{
		"catch": {
			description: "Catch a Pokemon and add it to your Pokedex",
			callback:    commands.CatchFunc(client, trainer),
		},
		"exit": {
			description: "Exit the Pokedex",
			callback:    commands.ExitProgram,
		},
		"explore": {
			description: "List all the Pokemon in a given area",
			callback:    commands.ExploreFunc(client, trainer),
		},
		"help": {
			description: "Display the help message",
			callback:    nil,
		},
		"inspect": {
			description: "Inspect a Pokemon from your Pokedex",
			callback:    commands.InspectFunc(trainer),
		},
		"map": {
			description: "Display the next 20 locations in the Pokemon world",
			callback:    commands.MapFunc(client, trainer),
		},
		"mapb": {
			description: "Display the previous 20 locations in the Pokemon world",
			callback:    commands.MapBFunc(client, trainer),
		},
		"pokedex": {
			description: "List the names of all the Pokemon in your Pokedex",
			callback:    commands.PokedexFunc(trainer),
		},
		"release": {
			description: "Release a Pokemon back into the wild",
			callback:    commands.ReleaseFunc(trainer),
		},
		"visit": {
			description: "Visit a location area",
			callback:    commands.VisitFunc(client, trainer),
		},
	}

	summaries := summaryMap(commandMap)

	commandMap["help"] = command{
		description: "Displays a help message",
		callback:    commands.HelpFunc(summaries),
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
