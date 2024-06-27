package app

import (
	"encoding/json"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	squareservices "github.com/longvu727/FootballSquaresLibs/services/square_microservices"
	"github.com/longvu727/FootballSquaresLibs/util"
)

type GetGameParams struct {
	GameGUID string `json:"game_guid"`
}
type GetGameResponse struct {
	Game            gameservices.Game                                      `json:"game"`
	Square          squareservices.Square                                  `json:"square"`
	FootballSquares []footballsquaregameservices.FootballSquareGameElement `json:"football_squares"`

	ErrorMessage string `json:"error_message"`
}

func (response GetGameResponse) ToJson() []byte {
	jsonStr, _ := json.Marshal(response)
	return jsonStr
}

func GetFootballSquareGame(getGameParams GetGameParams, config *util.Config) (*GetGameResponse, error) {
	var getGameResponse GetGameResponse

	getGameByGUID := gameservices.GetGameByGUIDService{GameGUID: getGameParams.GameGUID}
	getGameByGUIDResponse, err := getGameByGUID.Request(config)
	if err != nil {
		return &getGameResponse, nil
	}

	getFootballSquareGameByGameIDService := footballsquaregameservices.GetFootballSquareGameByGameIDService{
		GameID: int(getGameByGUIDResponse.GameID),
	}
	getFootballSquareGameByGameIDResponse, err := getFootballSquareGameByGameIDService.Request(config)
	if err != nil || len(getFootballSquareGameByGameIDResponse.FootballSquares) == 0 {
		return &getGameResponse, nil
	}

	getSquareService := squareservices.GetSquareService{
		SquareID: int(getFootballSquareGameByGameIDResponse.FootballSquares[0].SquareID),
	}
	getSquareResponse, err := getSquareService.Request(config)
	if err != nil {
		return &getGameResponse, nil
	}

	getGameResponse.FootballSquares = getFootballSquareGameByGameIDResponse.FootballSquares
	getGameResponse.Square = getSquareResponse.Square
	getGameResponse.Game = getGameByGUIDResponse.Game

	return &getGameResponse, nil
}
