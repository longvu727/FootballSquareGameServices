package websocketroutes

import (
	"encoding/json"
	"footballsquaregameservices/app"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type WebSocketRoutesInterface interface {
	Register(mux *http.ServeMux, resources *resources.Resources)
}

type WebSocketRoutes struct {
	Apps app.FootballSquareGame
}

func NewWebSocketRoutes(app app.FootballSquareGame) WebSocketRoutesInterface {
	return &WebSocketRoutes{
		Apps: app,
	}
}

func (routes *WebSocketRoutes) Register(mux *http.ServeMux, resources *resources.Resources) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	mux.HandleFunc("GET /Subscribe/GetGame/{game_guid}", func(w http.ResponseWriter, r *http.Request) {
		routes.SubscribeGame(w, r, resources, upgrader)
	})
}

func (routes *WebSocketRoutes) SubscribeGame(
	writer http.ResponseWriter,
	request *http.Request,
	resources *resources.Resources,
	upgrader websocket.Upgrader,
) {

	log.Printf("Received websocket request for %s\n", request.URL.Path)

	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer connection.Close()

	gameGUID := request.PathValue("game_guid")

	squareReservedChannel := resources.RedisClient.Subscribe(resources.Context, "SquareReserved:"+gameGUID)
	defer squareReservedChannel.Close()

	subscribeGameChannel := resources.RedisClient.Subscribe(resources.Context, "SubscribeGame:"+gameGUID)
	defer subscribeGameChannel.Close()

	response, err := routes.Apps.GetFootballSquareGame(app.GetGameParams{GameGUID: gameGUID}, resources)
	if err != nil {
		return
	}

	if err := connection.WriteJSON(response); err != nil {
		return
	}

	go func() {
		for {
			_, err := squareReservedChannel.ReceiveMessage(resources.Context)
			if err != nil {
				log.Println("SquareReservedChannel receive message error: " + err.Error())
			}

			response, err := routes.Apps.GetFootballSquareGame(app.GetGameParams{GameGUID: gameGUID}, resources)
			if err != nil {
				return
			}

			responseJSON, _ := json.Marshal(response)

			if err := resources.RedisClient.Publish(resources.Context, "SubscribeGame:"+gameGUID, responseJSON); err != nil {
				log.Println("SquareReservedChannel Publish error:", err)
			}
		}
	}()

	for msg := range subscribeGameChannel.Channel() {
		if err := connection.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			return
		}
	}
}
