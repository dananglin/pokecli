package main

import "maps"

func summaryMap(commandMap map[string]command) map[string]string {
	summaries := make(map[string]string)

	for key, value := range maps.All(commandMap) {
		summaries[key] = value.description
	}

	return summaries
}
