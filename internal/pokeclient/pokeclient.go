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
	baseURL string = "https://pokeapi.co/api/v2"

	LocationAreaPath = baseURL + "/location-area"
	PokemonPath      = baseURL + "/pokemon"
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

	dataFromCache, exists := c.cache.Get(url)
	if exists {
		fmt.Println("Using data from cache.")

		if err := decodeJSON(dataFromCache, &list); err != nil {
			return pokeapi.NamedAPIResourceList{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return list, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf("error creating the HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf("error getting the response from the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf(
			"received a bad status from %s: (%d) %s",
			url,
			resp.StatusCode,
			resp.Status,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf(
			"unable to read the response from the server: %w",
			err,
		)
	}

	if err := decodeJSON(body, &list); err != nil {
		return pokeapi.NamedAPIResourceList{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, body)

	return list, nil
}

func (c *Client) GetLocationArea(location string) (pokeapi.LocationArea, error) {
	var locationArea pokeapi.LocationArea

	url := LocationAreaPath + "/" + location + "/"

	dataFromCache, exists := c.cache.Get(url)
	if exists {
		fmt.Println("Using data from cache.")

		if err := decodeJSON(dataFromCache, &locationArea); err != nil {
			return pokeapi.LocationArea{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return locationArea, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf("error creating the HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf("error getting the response from the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return pokeapi.LocationArea{}, fmt.Errorf(
			"received a bad status from %s: (%d) %s",
			url,
			resp.StatusCode,
			resp.Status,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf(
			"unable to read the response from the server: %w",
			err,
		)
	}

	if err := decodeJSON(body, &locationArea); err != nil {
		return pokeapi.LocationArea{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, body)

	return locationArea, nil
}

func (c *Client) GetPokemon(pokemonName string) (pokeapi.Pokemon, error) {
	var pokemon pokeapi.Pokemon

	url := PokemonPath + "/" + pokemonName + "/"

	dataFromCache, exists := c.cache.Get(url)
	if exists {
		fmt.Println("Using data from cache.")

		if err := decodeJSON(dataFromCache, &pokemon); err != nil {
			return pokeapi.Pokemon{}, fmt.Errorf("unable to decode the data from the cache: %w", err)
		}

		return pokemon, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("error creating the HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("error getting the response from the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return pokeapi.Pokemon{}, fmt.Errorf(
			"received a bad status from %s: (%d) %s",
			url,
			resp.StatusCode,
			resp.Status,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf(
			"unable to read the response from the server: %w",
			err,
		)
	}

	if err := decodeJSON(body, &pokemon); err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("unable to decode the data from the server: %w", err)
	}

	c.cache.Add(url, body)

	return pokemon, nil
}

func decodeJSON(data []byte, value any) error {
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return nil
}
