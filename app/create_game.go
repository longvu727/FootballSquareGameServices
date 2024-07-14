package app

import (
	"encoding/json"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	squareservices "github.com/longvu727/FootballSquaresLibs/services/square_microservices"
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

	createSquareService := squareservices.CreateSquareService{
		SquareSize: int(createGameParams.SquareSize),
		Sport:      createGameParams.Sport,
	}

	createSquareServiceResponse, err := createSquareService.Request(&resources.Config)
	if err != nil {
		return &createGameResponse, nil
	}

	createGameService := gameservices.CreateGameService{
		Sport:      createGameParams.Sport,
		SquareSize: createGameParams.SquareSize,
		TeamA:      createGameParams.TeamA,
		TeamB:      createGameParams.TeamB,
	}

	createGameServiceResponse, err := createGameService.Request(&resources.Config)
	if err != nil {
		return &createGameResponse, nil
	}

	createFootballSquareGameService := footballsquaregameservices.CreateFootballSquareGameService{
		GameID:     int(createGameServiceResponse.GameID),
		SquareID:   createSquareServiceResponse.SquareID,
		SquareSize: int(createGameParams.SquareSize),
	}

	_, err = createFootballSquareGameService.Request(&resources.Config)
	if err != nil {
		return &createGameResponse, nil
	}

	createGameResponse.GameGUID = createGameServiceResponse.GameGUID

	return &createGameResponse, nil
}
