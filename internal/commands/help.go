package commands

import (
	"fmt"
	"maps"
	"slices"
)

func HelpFunc(summaries map[string]string) CommandFunc {
	return func(_ []string) error {
		keys := []string{}

		for key := range maps.All(summaries) {
			keys = append(keys, key)
		}

		slices.Sort(keys)

		fmt.Printf("\nCommands:\n")

		for _, key := range slices.All(keys) {
			fmt.Printf("\n%s: %s", key, summaries[key])
		}

		fmt.Printf("\n\n")

		return nil
	}
}
