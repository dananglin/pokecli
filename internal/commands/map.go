package commands

import (
	"fmt"
	"slices"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokeclient"
	"codeflow.dananglin.me.uk/apollo/pokedex/internal/poketrainer"
)

func MapFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) CommandFunc {
	return func(_ []string) error {
		url := trainer.NextLocationArea()
		if url == nil {
			url = new(string)
			*url = pokeclient.LocationAreaPath
		}

		return printResourceList(client, *url, trainer.UpdateLocationAreas)
	}
}

func MapBFunc(client *pokeclient.Client, trainer *poketrainer.Trainer) CommandFunc {
	return func(_ []string) error {
		url := trainer.PreviousLocationArea()
		if url == nil {
			return fmt.Errorf("no previous locations available")
		}

		return printResourceList(client, *url, trainer.UpdateLocationAreas)
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
