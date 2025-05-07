package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lucashthiele/pokedexcli/model"
)

func GetPokemons(url string) (model.PokemonResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return model.PokemonResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.PokemonResponse{}, err
	}

	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.PokemonResponse{}, fmt.Errorf("BodyContent: %v", body)
	}

	pokemon, err := UnmarshalPokemonResponse(body)
	if err != nil {
		return model.PokemonResponse{}, err
	}

	return pokemon, nil
}

func UnmarshalPokemonResponse(response []byte) (model.PokemonResponse, error) {
	pokemon := model.PokemonResponse{}

	err := json.Unmarshal(response, &pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, nil
}

func MarshalPokemon(pokemon model.PokemonResponse) ([]byte, error) {
	data, err := json.Marshal(pokemon)
	if err != nil {
		return nil, err
	}

	return data, nil
}
