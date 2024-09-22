package commands

import "codeflow.dananglin.me.uk/apollo/pokecli/internal/poketrainer"

func PokedexFunc(trainer *poketrainer.Trainer) CommandFunc {
	return func(_ []string) error {
		trainer.ListAllPokemonFromPokedex()

		return nil
	}
}
