package app

import (
	"encoding/json"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/services"
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

func CreateFootballSquareGame(request *http.Request) (*CreateGameResponse, error) {
	var createGameResponse CreateGameResponse

	var createGameParams CreateGameParams
	json.NewDecoder(request.Body).Decode(&createGameParams)

	createSquareService := services.CreateSquare{
		SquareSize: int(createGameParams.SquareSize),
		Sport:      createGameParams.Sport,
	}

	createSquareServiceResponse, err := createSquareService.Request()
	if err != nil {
		return &createGameResponse, nil
	}

	createGameService := services.CreateGame{
		Sport:      createGameParams.Sport,
		SquareSize: createGameParams.SquareSize,
		TeamA:      createGameParams.TeamA,
		TeamB:      createGameParams.TeamB,
	}

	createGameServiceResponse, err := createGameService.Request()
	if err != nil {
		return &createGameResponse, nil
	}

	createFootballSquareGameService := services.CreateFootballSquareGame{
		GameID:     int(createGameServiceResponse.GameID),
		SquareID:   createSquareServiceResponse.SquareID,
		SquareSize: int(createGameParams.SquareSize),
	}

	_, err = createFootballSquareGameService.Request()
	if err != nil {
		return &createGameResponse, nil
	}

	createGameResponse.GameGUID = createGameServiceResponse.GameGUID

	return &createGameResponse, nil
}
