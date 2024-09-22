package commands

import (
	"errors"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/pokecli/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokecli/internal/poketrainer"
)

func VisitFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) CommandFunc {
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
