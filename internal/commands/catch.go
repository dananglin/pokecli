package commands

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"

	"codeflow.dananglin.me.uk/apollo/pokecli/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokecli/internal/poketrainer"
)

func CatchFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) CommandFunc {
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

		if _, caught := trainer.GetPokemonFromPokedex(pokemonName); caught {
			return fmt.Errorf(
				"you've already caught a %s",
				pokemonName,
			)
		}

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

		if caught := success(chance); caught {
			trainer.AddPokemonToPokedex(pokemonName, pokemonDetails)
			fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", pokemonName)
		} else {
			fmt.Printf("%s escaped!\n", pokemonName)
		}

		return nil
	}
}

func success(chance int) bool {
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
