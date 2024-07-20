package routes

import (
	"footballsquaregameservices/app"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type WebSocketHandler = func(writer http.ResponseWriter, request *http.Request, resources *resources.Resources, upgrader websocket.Upgrader)

func (routes *Routes) registerWebSocketRoutes(mux *http.ServeMux, resources *resources.Resources) {
	routesHandlersMap := map[string]WebSocketHandler{
		"GET /Subscribe/Game/{game_guid}": routes.SubscribeGame,
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	for route, handler := range routesHandlersMap {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, resources, upgrader)
		})
	}

}

func (routes *Routes) SubscribeGame(writer http.ResponseWriter, request *http.Request, resources *resources.Resources, upgrader websocket.Upgrader) {
	log.Printf("Received websocket request for %s\n", request.URL.Path)

	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer connection.Close()

	getGameParams := app.GetGameParams{
		GameGUID: request.PathValue("game_guid"),
	}

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		response, err := routes.Apps.GetFootballSquareGame(getGameParams, resources)
		if err != nil {
			connection.Close()
		}

		connection.WriteJSON(response)
	}

}
