package pokecache_test

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"codeflow.dananglin.me.uk/apollo/pokedex/internal/pokecache"
)

const (
	keyNotFoundFormat = "The key %q was not found after adding it to the cache"
	keyFoundAfterCleanupFormat = "The key %q was found after cache cleanup"
)

func TestCacheAddGet(t *testing.T) {
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://example.org/path",
			value: []byte("testdata"),
		},
		{
			key:   "https://example.org/api/v1/path",
			value: []byte(`{"version": "v1.0.0", "key": "value"}`),
		},
	}

	interval := 1 * time.Minute

	cache := pokecache.NewCache(interval)

	testFunc := func(key string, value []byte) func(*testing.T) {
		return func(t *testing.T) {
			cache.Add(key, value)
			gotBytes, exists := cache.Get(key)

			if !exists {
				t.Fatalf(keyNotFoundFormat, key)
			}

			want := string(value)
			got := string(gotBytes)

			if got != want {
				t.Errorf("Unexpected value retrieved from the cache: want %s, got %s", want, got)
			}
		}
	}

	for ind, testcase := range slices.All(cases) {
		t.Run(fmt.Sprintf("Test case: %d", ind+1), testFunc(testcase.key, testcase.value))
	}
}

func TestReadLoop(t *testing.T) {
	const (
		baseTime = 5 * time.Millisecond
		waitTime = 10 * baseTime
	)

	key := "https://example.org/api/v1/path"
	value := []byte(`{"version": "v1.0.0", "key": "value"}`)

	cache := pokecache.NewCache(baseTime)

	cache.Add(key, value)

	if _, exists := cache.Get(key); !exists {
		t.Fatalf(keyNotFoundFormat, key)
	}

	time.Sleep(waitTime)

	if _, exists := cache.Get(key); exists {
		t.Errorf(keyFoundAfterCleanupFormat, key)
	}
}
