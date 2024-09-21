package commands

import (
	"errors"
	"fmt"
	"slices"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/poketrainer"
)

func InspectFunc(trainer *poketrainer.Trainer) CommandFunc {
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
