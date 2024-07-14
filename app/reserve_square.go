package app

import (
	"encoding/json"

	footballsquaregameservices "github.com/longvu727/FootballSquaresLibs/services/football_square_game_microservices"
	gameservices "github.com/longvu727/FootballSquaresLibs/services/game_microservices"
	userservices "github.com/longvu727/FootballSquaresLibs/services/user_microservices"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type ReserveSquareParams struct {
	GameGUID    string `json:"game_guid"`
	UserGUID    string `json:"user_guid"`
	RowIndex    int    `json:"row_index"`
	ColumnIndex int    `json:"column_index"`
}
type ReserveSquareResponse struct {
	Reserved     bool   `json:"reserved"`
	ErrorMessage string `json:"error_message"`
}

func (response ReserveSquareResponse) ToJson() []byte {
	jsonStr, _ := json.Marshal(response)
	return jsonStr
}

func (footballSquareGameApp *FootballSquareGameApp) ReserveSquare(reserveSquareParams ReserveSquareParams, resources *resources.Resources) (*ReserveSquareResponse, error) {
	reserveSquareResponse := &ReserveSquareResponse{Reserved: false}

	getGameByGUID := gameservices.GetGameByGUIDService{GameGUID: reserveSquareParams.GameGUID}
	getGameByGUIDResponse, err := getGameByGUID.Request(&resources.Config)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to find game"
		return reserveSquareResponse, err
	}

	getUserByGUID := userservices.GetUserByGUIDService{UserGUID: reserveSquareParams.UserGUID}
	getUserByGUIDResponse, err := getUserByGUID.Request(&resources.Config)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to find user"
		return reserveSquareResponse, err
	}

	reserveSquare := footballsquaregameservices.ReserveFootballSquareService{
		GameID:      int(getGameByGUIDResponse.GameID),
		UserID:      int(getUserByGUIDResponse.UserID),
		RowIndex:    reserveSquareParams.RowIndex,
		ColumnIndex: reserveSquareParams.ColumnIndex,
	}

	_, err = reserveSquare.Request(&resources.Config)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to reserve square"
		return reserveSquareResponse, err
	}

	reserveSquareResponse.Reserved = true
	return reserveSquareResponse, nil
}
