package app

import (
	"encoding/json"

	"github.com/longvu727/FootballSquaresLibs/services"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type CreateGameParams struct {
	Sport      string `json:"sport"`
	SquareSize int32  `json:"square_size"`
	TeamA      string `json:"team_a"`
	TeamB      string `json:"team_b"`
}
type CreateGameResponse struct {
	GameGUID     string `json:"game_guid"`
	ErrorMessage string `json:"error_message"`
}

func (response CreateGameResponse) ToJson() []byte {
	jsonStr, _ := json.Marshal(response)
	return jsonStr
}

func (footballSquareGameApp *FootballSquareGameApp) CreateFootballSquareGame(createGameParams CreateGameParams, resources *resources.Resources) (*CreateGameResponse, error) {
	var createGameResponse CreateGameResponse

	createSquareRequest := services.CreateSquareRequest{
		SquareSize: int(createGameParams.SquareSize),
		Sport:      createGameParams.Sport,
	}

	createSquareServiceResponse, err := resources.Services.CreateSquare(&resources.Config, createSquareRequest)
	if err != nil {
		return &createGameResponse, nil
	}

	createGameRequest := services.CreateGameRequest{
		Sport:      createGameParams.Sport,
		SquareSize: createGameParams.SquareSize,
		TeamA:      createGameParams.TeamA,
		TeamB:      createGameParams.TeamB,
	}

	createGameServiceResponse, err := resources.Services.CreateGame(&resources.Config, createGameRequest)
	if err != nil {
		return &createGameResponse, nil
	}

	createFootballSquareGameRequest := services.CreateFootballSquareGameRequest{
		GameID:     int(createGameServiceResponse.GameID),
		SquareID:   createSquareServiceResponse.SquareID,
		SquareSize: int(createGameParams.SquareSize),
	}

	_, err = resources.Services.CreateFootballSquareGame(&resources.Config, createFootballSquareGameRequest)
	if err != nil {
		return &createGameResponse, nil
	}

	createGameResponse.GameGUID = createGameServiceResponse.GameGUID

	return &createGameResponse, nil
}
