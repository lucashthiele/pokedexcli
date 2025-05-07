package model

import "fmt"

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
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
type Pokemon struct {
	Name           string  `json:"name,omitempty"`
	BaseExperience int     `json:"base_experience,omitempty"`
	Height         int     `json:"height,omitempty"`
	Weight         int     `json:"weight,omitempty"`
	Stats          []Stats `json:"stats,omitempty"`
	Types          []Types `json:"types,omitempty"`
}

func (p *Pokemon) ToString() string {
	result := fmt.Sprintf("Name: %s\n", p.Name)
	result += fmt.Sprintf("Height: %d\n", p.Height)
	result += fmt.Sprintf("Weight: %d\n", p.Weight)
	result += fmt.Sprintf("Stats:\n")

	for _, pokemonStat := range p.Stats {
		result += fmt.Sprintf(" -%s: %d\n", pokemonStat.Stat.Name, pokemonStat.BaseStat)
	}

	result += fmt.Sprintf("Types:\n")

	for _, pokemonType := range p.Types {
		result += fmt.Sprintf(" - %s\n", pokemonType.Type.Name)
	}

	return result
}
