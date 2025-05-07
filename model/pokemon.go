package model

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}
