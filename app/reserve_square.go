package app

import (
	"encoding/json"

	"github.com/longvu727/FootballSquaresLibs/services"
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

	getGameByGUIDRequest := services.GetGameByGUIDRequest{GameGUID: reserveSquareParams.GameGUID}
	getGameByGUIDResponse, err := resources.Services.GetGameByGUID(&resources.Config, getGameByGUIDRequest)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to find game"
		return reserveSquareResponse, err
	}

	getUserByGUIDRequest := services.GetUserByGUIDRequest{UserGUID: reserveSquareParams.UserGUID}
	getUserByGUIDResponse, err := resources.Services.GetUserByGUID(&resources.Config, getUserByGUIDRequest)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to find user"
		return reserveSquareResponse, err
	}

	reserveSquareRequest := services.ReserveFootballSquareRequest{
		GameID:      int(getGameByGUIDResponse.GameID),
		UserID:      int(getUserByGUIDResponse.UserID),
		RowIndex:    reserveSquareParams.RowIndex,
		ColumnIndex: reserveSquareParams.ColumnIndex,
	}

	_, err = resources.Services.ReserveFootballSquare(&resources.Config, reserveSquareRequest)
	if err != nil {
		reserveSquareResponse.ErrorMessage = "unable to reserve square"
		return reserveSquareResponse, err
	}

	reserveSquareResponse.Reserved = true
	return reserveSquareResponse, nil
}
