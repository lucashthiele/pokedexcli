package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lucashthiele/pokedexcli/model"
)

func GetLocations(url string) (model.LocationResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return model.LocationResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.LocationResponse{}, err
	}

	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.LocationResponse{}, fmt.Errorf(fmt.Sprintf("BodyContent: %v", body))
	}

	location := model.LocationResponse{}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return model.LocationResponse{}, err
	}

	return location, nil
}
