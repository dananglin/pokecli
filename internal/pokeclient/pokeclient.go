package pokeclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/api/pokeapi"
	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokecache"
)

const (
	baseURL string = "https://pokeapi.co"

	LocationAreaPath = baseURL + "/api/v2/location-area"
	PokemonPath      = baseURL + "/api/v2/pokemon"
)

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
	timeout    time.Duration
}

func NewClient(cacheCleanupInterval, timeout time.Duration) *Client {
	cache := pokecache.NewCache(cacheCleanupInterval)

	client := Client{
		httpClient: http.Client{},
		cache:      cache,
		timeout:    timeout,
	}

	return &client
}

func (c *Client) GetNamedAPIResourceList(url string) (pokeapi.NamedAPIResourceList, error) {
	var list pokeapi.NamedAPIResourceList

	data, exists := c.cache.Get(url)
	if exists {
		fmt.Println("(using data from cache)")

		if err := decodeJSON(data, &list); err != nil {
			return pokeapi.NamedAPIResourceList{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return list, nil
	}

	data, err := c.sendRequest(url)
	if err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf(
			"received an error after sending the request to the server: %w",
			err,
		)
	}

	if err := decodeJSON(data, &list); err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, data)

	return list, nil
}

func (c *Client) GetLocationArea(location string) (pokeapi.LocationArea, error) {
	var locationArea pokeapi.LocationArea

	url := LocationAreaPath + "/" + location + "/"

	data, exists := c.cache.Get(url)
	if exists {
		fmt.Println("(using data from cache)")

		if err := decodeJSON(data, &locationArea); err != nil {
			return pokeapi.LocationArea{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return locationArea, nil
	}

	data, err := c.sendRequest(url)
	if err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf(
			"received an error after sending the request to the server: %w",
			err,
		)
	}

	if err := decodeJSON(data, &locationArea); err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, data)

	return locationArea, nil
}

func (c *Client) GetPokemon(pokemonName string) (pokeapi.Pokemon, error) {
	var pokemon pokeapi.Pokemon

	url := PokemonPath + "/" + pokemonName + "/"

	data, exists := c.cache.Get(url)
	if exists {
		fmt.Println("(using data from cache)")

		if err := decodeJSON(data, &pokemon); err != nil {
			return pokeapi.Pokemon{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return pokemon, nil
	}

	data, err := c.sendRequest(url)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf(
			"received an error after sending the request to the server: %w",
			err,
		)
	}

	if err := decodeJSON(data, &pokemon); err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, data)

	return pokemon, nil
}

func (c *Client) GetPokemonLocationAreas(url string) ([]pokeapi.LocationAreaEncounter, error) {
	var locationAreaEncounters []pokeapi.LocationAreaEncounter

	data, exists := c.cache.Get(url)
	if exists {
		fmt.Println("(using data from cache)")

		if err := decodeJSON(data, &locationAreaEncounters); err != nil {
			return []pokeapi.LocationAreaEncounter{}, fmt.Errorf(
				"unable to decode the data from the cache: %w",
				err,
			)
		}
	}

	data, err := c.sendRequest(url)
	if err != nil {
		return []pokeapi.LocationAreaEncounter{}, fmt.Errorf(
			"received an error after sending the request to the server: %w",
			err,
		)
	}

	if err := decodeJSON(data, &locationAreaEncounters); err != nil {
		return []pokeapi.LocationAreaEncounter{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	return locationAreaEncounters, nil
}

func (c *Client) sendRequest(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("error creating the HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("error getting the response from the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return []byte{}, fmt.Errorf(
			"received a bad status from %s: (%d) %s",
			url,
			resp.StatusCode,
			resp.Status,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf(
			"unable to read the response from the server: %w",
			err,
		)
	}

	return data, nil
}

func decodeJSON(data []byte, value any) error {
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return nil
}
