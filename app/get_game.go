package app

import (
	"encoding/json"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	squareservices "github.com/longvu727/FootballSquaresLibs/services/square_microservices"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type GetGameParams struct {
	GameGUID string `json:"game_guid"`
}

type FootballSquare struct {
	ColumnIndex        int    `json:"column_index"`
	RowIndex           int    `json:"row_index"`
	WinnerQuaterNumber int    `json:"winner_quater_number"`
	Winner             bool   `json:"winner"`
	UserGUID           string `json:"user_guid"`
	UserAlias          string `json:"user_alias"`
	UserName           string `json:"user_name"`
}

type GetGameResponse struct {
	//Game
	GameGUID string `json:"game_guid"`
	Sport    string `json:"sport"`
	TeamA    string `json:"team_a"`
	TeamB    string `json:"team_b"`

	//Square
	SquareSize   int    `json:"square_size"`
	RowPoints    string `json:"row_points"`
	ColumnPoints string `json:"column_points"`

	FootballSquares []FootballSquare `json:"football_squares"`

	ErrorMessage string `json:"error_message"`
}

func (response GetGameResponse) ToJson() []byte {
	jsonStr, _ := json.Marshal(response)
	return jsonStr
}

func (game *GameApp) GetFootballSquareGame(getGameParams GetGameParams, resources *resources.Resources) (*GetGameResponse, error) {
	var getGameResponse GetGameResponse

	getGameByGUID := gameservices.GetGameByGUIDService{GameGUID: getGameParams.GameGUID}
	getGameByGUIDResponse, err := getGameByGUID.Request(&resources.Config)
	if err != nil {
		return &getGameResponse, nil
	}

	getFootballSquareGameByGameIDService := footballsquaregameservices.GetFootballSquareGameByGameIDService{
		GameID: int(getGameByGUIDResponse.GameID),
	}
	getFootballSquareGameByGameIDResponse, err := getFootballSquareGameByGameIDService.Request(&resources.Config)
	if err != nil || len(getFootballSquareGameByGameIDResponse.FootballSquares) == 0 {
		return &getGameResponse, nil
	}

	getSquareService := squareservices.GetSquareService{
		SquareID: int(getFootballSquareGameByGameIDResponse.FootballSquares[0].SquareID),
	}
	getSquareResponse, err := getSquareService.Request(&resources.Config)
	if err != nil {
		return &getGameResponse, nil
	}

	getGameResponse.GameGUID = getGameByGUIDResponse.GameGUID
	getGameResponse.Sport = getGameByGUIDResponse.Sport
	getGameResponse.TeamA = getGameByGUIDResponse.TeamA
	getGameResponse.TeamB = getGameByGUIDResponse.TeamB

	getGameResponse.SquareSize = getSquareResponse.SquareSize
	getGameResponse.RowPoints = getSquareResponse.RowPoints
	getGameResponse.ColumnPoints = getSquareResponse.ColumnPoints

	for _, footballSquare := range getFootballSquareGameByGameIDResponse.FootballSquares {
		getGameResponse.FootballSquares = append(getGameResponse.FootballSquares, FootballSquare{
			ColumnIndex:        footballSquare.ColumnIndex,
			RowIndex:           footballSquare.RowIndex,
			WinnerQuaterNumber: footballSquare.WinnerQuaterNumber,
			Winner:             footballSquare.Winner,
			UserGUID:           "",
			UserAlias:          "",
			UserName:           "",
		})
	}

	return &getGameResponse, nil
}
