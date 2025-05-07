package model

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string  `json:"name,omitempty"`
	BaseExperience int     `json:"base_experience,omitempty"`
	Height         int     `json:"height,omitempty"`
	Weight         int     `json:"weight,omitempty"`
	Stats          []Stats `json:"stats,omitempty"`
	Types          []Types `json:"types,omitempty"`
}
type Stat struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Stats struct {
	BaseStat int  `json:"base_stat,omitempty"`
	Effort   int  `json:"effort,omitempty"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Types struct {
	Slot int  `json:"slot,omitempty"`
	Type Type `json:"type"`
}
