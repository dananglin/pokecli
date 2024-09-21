package commands

import (
	"fmt"
	"slices"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokedex/internal/poketrainer"
)

func ExploreFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) CommandFunc {
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
