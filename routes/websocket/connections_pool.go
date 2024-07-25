package websocketroutes

import (
	"footballsquaregameservices/app"

	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
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
	resources     *resources.Resources
}

func newSocketConnectionsPool(resources *resources.Resources) *socketConnectionsPool {
	socketConnectionsPool := &socketConnectionsPool{
		SubscribeGame: subscribeGameConnections{
			gameGUIDConnections: make(map[string]gameConnections),
			broadcast:           make(chan subscribeGameBroadcastData),
		},
		resources: resources,
	}

	return socketConnectionsPool
}

/*func (routes *HTTPRoutes) cacheGetGameResponse(gameGUID string, resources *resources.Resources, responseStr string) error {
	redisKey := routes.getGetGameCacheKey(gameGUID)
	return resources.RedisClient.Set(resources.Context, redisKey, responseStr, time.Hour).Err()
}
*/
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
