package routes

import (
	"footballsquaregameservices/app"
	"log"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type RoutesInterface interface {
	Register(resources *resources.Resources) *http.ServeMux
}

type Routes struct {
	Apps app.FootballSquareGame
}

func NewRoutes() RoutesInterface {
	return &Routes{
		Apps: app.NewFootballSquareGameApp(),
	}
}

func (routes *Routes) Register(resources *resources.Resources) *http.ServeMux {
	log.Println("Registering routes")
	mux := http.NewServeMux()

	routes.registerHttpRoutes(mux, resources)
	routes.registerWebSocketRoutes(mux, resources)

	return mux
}
