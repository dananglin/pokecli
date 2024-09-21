package commands

import (
	"errors"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/poketrainer"
)

func ReleaseFunc(trainer *poketrainer.Trainer) CommandFunc {
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

		if _, caught := trainer.GetPokemonFromPokedex(pokemonName); !caught {
			return fmt.Errorf(
				"you haven't caught a %s",
				pokemonName,
			)
		}

		trainer.RemovePokemonFromPokedex(pokemonName)

		fmt.Printf("%s was released back into the wild.\n", pokemonName)

		return nil
	}
}
