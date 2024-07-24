package websocketroutes

import (
	"encoding/json"
	"footballsquaregameservices/app"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/redis/go-redis/v9"
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

	socketConnectionsPool := newSocketConnectionsPool()

	mux.HandleFunc("GET /Subscribe/GetGame/{game_guid}", func(w http.ResponseWriter, r *http.Request) {
		routes.SubscribeGame(w, r, resources, upgrader, socketConnectionsPool)
	})
}

func (routes *WebSocketRoutes) SubscribeGame(
	writer http.ResponseWriter,
	request *http.Request,
	resources *resources.Resources,
	upgrader websocket.Upgrader,
	connections *socketConnectionsPool,
) {

	log.Printf("Received websocket request for %s\n", request.URL.Path)

	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer connection.Close()

	gameGUID := request.PathValue("game_guid")
	getGameParams := app.GetGameParams{
		GameGUID: gameGUID,
	}

	if _, ok := connections.getConnectionsByGameGUID(gameGUID); !ok {
		connections.newGameGUIDConnections(gameGUID)
	}
	connections.addGameGUIDConnection(gameGUID, connection)

	subscribeGameBroadcastData := subscribeGameBroadcastData{
		gameGUID: gameGUID,
	}

	oldTime := time.Now()
	newTime := time.Now()

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {

		log.Println("Ticked, GameGUID " + gameGUID)

		response, err := routes.Apps.GetFootballSquareGame(getGameParams, resources)
		if err != nil {
			connection.Close()
		}

		jsonStr, _ := json.Marshal(response)

		log.Println(string(jsonStr))

		subscribeGameBroadcastData.getFootballSquareGameResponse = *response

		connections.SubscribeGame.broadcast <- subscribeGameBroadcastData

		for newTime.Equal(oldTime) {
			tempTime := routes.getSquareReservedTime(resources, gameGUID)
			if tempTime != nil {
				newTime = *tempTime
			}
		}
		oldTime = newTime
	}

}

func (routes *WebSocketRoutes) getSquareReservedTime(resources *resources.Resources, gameGUID string) *time.Time {
	redisKey := "SquareReserved:" + gameGUID

	cachedResponseStr, err := resources.RedisClient.Get(resources.Context, redisKey).Result()

	if err == redis.Nil || err != nil {
		return nil
	}

	squareReservedTime, _ := time.Parse(time.UnixDate, cachedResponseStr)

	return &squareReservedTime

}
