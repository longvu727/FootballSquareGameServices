package app

import (
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type FootballSquareGame interface {
	CreateFootballSquareGame(createGameParams CreateGameParams, resources *resources.Resources) (*CreateGameResponse, error)
	GetFootballSquareGame(getGameParams GetGameParams, resources *resources.Resources) (*GetGameResponse, error)
	ReserveSquare(reserveSquareParams ReserveSquareParams, resources *resources.Resources) (*ReserveSquareResponse, error)
}

type FootballSquareGameApp struct{}

func NewFootballSquareGameApp() FootballSquareGame {
	return &FootballSquareGameApp{}
}
