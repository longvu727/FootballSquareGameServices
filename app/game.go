package app

import "github.com/longvu727/FootballSquaresLibs/util/resources"

type Game interface {
	CreateFootballSquareGame(createGameParams CreateGameParams, resources *resources.Resources) (*CreateGameResponse, error)
	GetFootballSquareGame(getGameParams GetGameParams, resources *resources.Resources) (*GetGameResponse, error)
}

type GameApp struct{}

func NewGameApp() Game {
	return &GameApp{}
}
