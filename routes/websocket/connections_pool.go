package websocketroutes

import (
	"footballsquaregameservices/app"
	"log"

	"github.com/gorilla/websocket"
)

type subscribeGameBroadcastData struct {
	gameGUID                      string
	getFootballSquareGameResponse app.GetGameResponse
}

type gameConnections map[*websocket.Conn]bool

type subscribeGameConnections struct {
	gameGUIDConnections map[string]gameConnections
	broadcast           chan subscribeGameBroadcastData
}

type socketConnectionsPool struct {
	SubscribeGame subscribeGameConnections
}

func newSocketConnectionsPool() *socketConnectionsPool {
	socketConnectionsPool := &socketConnectionsPool{
		SubscribeGame: subscribeGameConnections{
			gameGUIDConnections: make(map[string]gameConnections),
			broadcast:           make(chan subscribeGameBroadcastData),
		},
	}

	go socketConnectionsPool.broadcastGame()
	return socketConnectionsPool
}

func (connections *socketConnectionsPool) getConnectionsByGameGUID(gameGUID string) (gameConnections, bool) {
	gameConnections, ok := connections.SubscribeGame.gameGUIDConnections[gameGUID]

	return gameConnections, ok
}

func (connections *socketConnectionsPool) newGameGUIDConnections(gameGUID string) {
	connections.SubscribeGame.gameGUIDConnections[gameGUID] = make(gameConnections)
}

func (connections *socketConnectionsPool) addGameGUIDConnection(gameGUID string, connection *websocket.Conn) {
	connections.SubscribeGame.gameGUIDConnections[gameGUID][connection] = true
}

func (connections *socketConnectionsPool) broadcastGame() {
	for {
		subscribeGameBroadcastData := <-connections.SubscribeGame.broadcast
		log.Println("Broadcasted")

		for connection := range connections.SubscribeGame.gameGUIDConnections[subscribeGameBroadcastData.gameGUID] {
			err := connection.WriteJSON(subscribeGameBroadcastData.getFootballSquareGameResponse)
			if err != nil {
				log.Println(err)
				connection.Close()
				delete(connections.SubscribeGame.gameGUIDConnections[subscribeGameBroadcastData.gameGUID], connection)
			}
		}
	}
}
