package app

import (
	"encoding/json"
	"log"

	"github.com/longvu727/FootballSquaresLibs/services"
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

func (footballSquareGameApp *FootballSquareGameApp) GetFootballSquareGame(getGameParams GetGameParams, resources *resources.Resources) (*GetGameResponse, error) {
	var getGameResponse GetGameResponse

	getGameByGUIDRequest := services.GetGameByGUIDRequest{GameGUID: getGameParams.GameGUID}
	getGameByGUIDResponse, err := resources.Services.GetGameByGUID(&resources.Config, getGameByGUIDRequest)
	if err != nil {
		return &getGameResponse, nil
	}

	getFootballSquareGameByGameIDRequest := services.GetFootballSquareGameByGameIDRequest{
		GameID: int(getGameByGUIDResponse.GameID),
	}
	getFootballSquareGameByGameIDResponse, err := resources.Services.GetFootballSquareGameByGameID(&resources.Config, getFootballSquareGameByGameIDRequest)
	if err != nil || len(getFootballSquareGameByGameIDResponse.FootballSquares) == 0 {
		return &getGameResponse, nil
	}

	getSquareRequest := services.GetSquareRequest{
		SquareID: int(getFootballSquareGameByGameIDResponse.FootballSquares[0].SquareID),
	}
	getSquareResponse, err := resources.Services.GetSquare(&resources.Config, getSquareRequest)
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

		square := FootballSquare{
			ColumnIndex:        footballSquare.ColumnIndex,
			RowIndex:           footballSquare.RowIndex,
			WinnerQuaterNumber: footballSquare.WinnerQuaterNumber,
			Winner:             footballSquare.Winner,
		}

		if footballSquare.UserID > 0 {
			getUserRequest := services.GetUserRequest{
				UserID: int(footballSquare.UserID),
			}
			getUserResponse, err := resources.Services.GetUser(&resources.Config, getUserRequest)
			if err != nil {
				log.Printf("unable to find user, user_id: %d, error: %s", footballSquare.UserID, err.Error())
			} else {
				square.UserGUID = getUserResponse.UserGUID
				square.UserAlias = getUserResponse.Alias
				square.UserName = getUserResponse.UserName
			}
		}

		getGameResponse.FootballSquares = append(getGameResponse.FootballSquares, square)
	}

	return &getGameResponse, nil
}
