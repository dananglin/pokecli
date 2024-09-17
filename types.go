package main

// LocationArea is a section of areas such as floors in a building or a cave.
type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             NamedAPIResource      `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource        `json:"encounter_method"`
	VersionDetails  EncounterVersionDetails `json:"version_details"`
}

type EncounterVersionDetails struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

// PokemonEncounter is details of a possible Pokemon encounter.
type PokemonEncounter struct {
	Pokemon        NamedAPIResource          `json:"pokemon"`
	VersionDetails []VersionEncounterDetails `json:"version_details"`
}

type VersionEncounterDetails struct {
	Version          NamedAPIResource `json:"version"`
	MaxChance        int              `json:"max_chance"`
	EncounterDetails []Encounter      `json:"encounter_details"`
}

type Encounter struct {
	MinLevel        int                `json:"min_level"`
	MaxLevel        int                `json:"max_level"`
	ConditionValues []NamedAPIResource `json:"condition_values"`
	Chance          int                `json:"chance"`
	Method          NamedAPIResource   `json:"method"`
}

type NamedAPIResourceList struct {
	Count    int                `json:"count"`
	Next     *string             `json:"next"`
	Previous *string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
