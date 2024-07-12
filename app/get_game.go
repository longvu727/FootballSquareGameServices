package app

import (
	"encoding/json"
	"log"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	squareservices "github.com/longvu727/FootballSquaresLibs/services/square_microservices"
	userservices "github.com/longvu727/FootballSquaresLibs/services/user_microservices"
	"github.com/longvu727/FootballSquaresLibs/util"
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
			getUserService := userservices.GetUserService{
				UserID: int(footballSquare.UserID),
			}
			getUserResponse, err := getUserService.Request(config)
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
