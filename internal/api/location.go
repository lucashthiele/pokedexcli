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
		return model.LocationResponse{}, fmt.Errorf("BodyContent: %v", body)
	}

	location, err := UnmarshalLocationResponse(body)

	if err != nil {
		return model.LocationResponse{}, err
	}

	return location, nil
}

func UnmarshalLocationResponse(response []byte) (model.LocationResponse, error) {
	location := model.LocationResponse{}

	err := json.Unmarshal(response, &location)
	if err != nil {
		return location, err
	}

	return location, nil
}

func MarshalLocation(location model.LocationResponse) ([]byte, error) {
	data, err := json.Marshal(location)
	if err != nil {
		return nil, err
	}

	return data, nil
}
