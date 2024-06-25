package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type CreateFootballSquareGameResponse struct {
	GameGUID     string `json:"game_guid"`
	ErrorMessage string `json:"error_message"`
}

type CreateFootballSquareGame struct {
	SquareSize int    `json:"square_size"`
	Sport      string `json:"sport"`
	Response   CreateFootballSquareGameResponse
}

func (service CreateFootballSquareGame) Request() (CreateFootballSquareGameResponse, error) {
	createFootballSquareGameResponse := CreateFootballSquareGameResponse{}

	client := &http.Client{}
	serviceJson, _ := json.Marshal(service)

	request, err := http.NewRequest("POST", "http://footballsquaregamemicroservices:3001/CreateFootballSquareGame", bytes.NewBuffer(serviceJson))
	if err != nil {
		return createFootballSquareGameResponse, err
	}

	response, err := client.Do(request)
	if err != nil {
		return createFootballSquareGameResponse, err
	}

	if response.StatusCode != http.StatusOK {
		errStr := `unable to create square`
		createFootballSquareGameResponse.ErrorMessage = errStr

		return createFootballSquareGameResponse, errors.New(errStr)
	}

	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(&createFootballSquareGameResponse)

	return createFootballSquareGameResponse, nil
}
