package app

import (
	"encoding/json"
	"net/http"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	squareservices "github.com/longvu727/FootballSquaresLibs/services/square_microservices"
	"github.com/longvu727/FootballSquaresLibs/util"
)

type GetGameParams struct {
	GameGUID string `json:"game_guid"`
}
type GetGameResponse struct {
	Game           gameservices.GetGameByGUIDResponse                               `json:"game"`
	Square         squareservices.GetSquareResponse                                 `json:"square"`
	FootballSquare footballsquaregameservices.GetFootballSquareGameByGameIDResponse `json:"football_square"`

	ErrorMessage string `json:"error_message"`
}

func (response GetGameResponse) ToJson() []byte {
	jsonStr, _ := json.Marshal(response)
	return jsonStr
}

func GetFootballSquareGame(request *http.Request, config *util.Config) (*GetGameResponse, error) {
	var getGameResponse GetGameResponse

	var getGameParams GetGameParams
	json.NewDecoder(request.Body).Decode(&getGameParams)

	getGameByGUID := gameservices.GetGameByGUID{GameGUID: getGameParams.GameGUID}
	getGameByGUIDResponse, err := getGameByGUID.Request(config)
	if err != nil {
		return &getGameResponse, nil
	}

	getFootballSquareGameByGameIDService := footballsquaregameservices.GetFootballSquareGameByGameID{
		GameID: int(getGameByGUIDResponse.GameID),
	}
	getFootballSquareGameByGameIDResponse, err := getFootballSquareGameByGameIDService.Request(config)
	if err != nil || len(getFootballSquareGameByGameIDResponse.FootballSquare) == 0 {
		return &getGameResponse, nil
	}

	getSquareService := squareservices.GetSquare{
		SquareID: int(getFootballSquareGameByGameIDResponse.FootballSquare[0].SquareID),
	}
	getSquareResponse, err := getSquareService.Request(config)
	if err != nil {
		return &getGameResponse, nil
	}

	getGameResponse.FootballSquare = getFootballSquareGameByGameIDResponse
	getGameResponse.Square = getSquareResponse
	getGameResponse.Game = getGameByGUIDResponse

	return &getGameResponse, nil
}
