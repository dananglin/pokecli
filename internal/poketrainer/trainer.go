package poketrainer

import (
	"fmt"
	"maps"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/api/pokeapi"
)

type Trainer struct {
	previousLocationArea    *string
	nextLocationArea        *string
	currentLocationAreaName string
	pokedex                 map[string]pokeapi.Pokemon
}

func NewTrainer() *Trainer {
	trainer := Trainer{
		previousLocationArea:    nil,
		nextLocationArea:        nil,
		currentLocationAreaName: "",
		pokedex:                 make(map[string]pokeapi.Pokemon),
	}

	return &trainer
}

func (t *Trainer) UpdateLocationAreas(previous, next *string) {
	t.previousLocationArea = previous
	t.nextLocationArea = next
}

func (t *Trainer) PreviousLocationArea() *string {
	return t.previousLocationArea
}

func (t *Trainer) NextLocationArea() *string {
	return t.nextLocationArea
}

func (t *Trainer) AddPokemonToPokedex(name string, details pokeapi.Pokemon) {
	t.pokedex[name] = details
}

func (t *Trainer) GetPokemonFromPokedex(name string) (pokeapi.Pokemon, bool) {
	details, ok := t.pokedex[name]

	return details, ok
}

func (t *Trainer) RemovePokemonFromPokedex(name string) {
	delete(t.pokedex, name)
}

func (t *Trainer) ListAllPokemonFromPokedex() {
	if len(t.pokedex) == 0 {
		fmt.Println("You have no Pokemon in your Pokedex.")

		return
	}

	fmt.Println("Your Pokedex:")

	for name := range maps.All(t.pokedex) {
		fmt.Println("  -", name)
	}
}

func (t *Trainer) CurrentLocationAreaName() string {
	return t.currentLocationAreaName
}

func (t *Trainer) UpdateCurrentLocationAreaName(locationName string) {
	t.currentLocationAreaName = locationName
}
