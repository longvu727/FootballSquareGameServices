package routes

import (
	"footballsquaregameservices/app"
	httproutes "footballsquaregameservices/routes/http"
	websocketroutes "footballsquaregameservices/routes/websocket"
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

	httpRoutes := httproutes.NewHTTPRoutes(routes.Apps)
	httpRoutes.Register(mux, resources)

	websocketRoutes := websocketroutes.NewWebSocketRoutes(routes.Apps)
	websocketRoutes.Register(mux, resources)

	return mux
}
