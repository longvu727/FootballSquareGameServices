package httproutes

import (
	"encoding/json"
	"fmt"
	"footballsquaregameservices/app"
	"log"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type HttpHandler = func(writer http.ResponseWriter, request *http.Request, resources *resources.Resources)

type HTTPRoutesInterface interface {
	Register(mux *http.ServeMux, resources *resources.Resources)
}

type HTTPRoutes struct {
	Apps app.FootballSquareGame
}

func NewHTTPRoutes(app app.FootballSquareGame) HTTPRoutesInterface {
	return &HTTPRoutes{
		Apps: app,
	}
}

func (routes *HTTPRoutes) Register(mux *http.ServeMux, resources *resources.Resources) {
	routesHandlersMap := map[string]HttpHandler{
		"/":                        routes.home,
		"POST /CreateGame":         routes.createGame,
		"GET /GetGame/{game_guid}": routes.getGame,
		"/GetEmptySquares":         routes.getEmptySquares,
		"POST /ReserveSquares":     routes.reserveSquares,
		"/SaveSquares":             routes.saveSquares,
		"/DeleteSquare":            routes.deleteSquare,
		"/GenerateNumber":          routes.generateNumber,
	}

	for route, handler := range routesHandlersMap {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", "*")
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
			} else {
				handler(w, r, resources)
			}
		})
	}
}

func (routes *HTTPRoutes) home(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	fmt.Fprintf(writer, "{\"Acknowledged\": true}")
}

func (routes *HTTPRoutes) createGame(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)

	var createGameParams app.CreateGameParams
	json.NewDecoder(request.Body).Decode(&createGameParams)

	response, err := routes.Apps.CreateFootballSquareGame(createGameParams, resources)

	if err != nil {
		response.ErrorMessage = `Unable to create game`
		responseStr, _ := json.Marshal(response)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(responseStr)

		return
	}

	responseStr, _ := json.Marshal(response)

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseStr)
}

func (routes *HTTPRoutes) getGame(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)

	getGameParams := app.GetGameParams{
		GameGUID: request.PathValue("game_guid"),
	}

	response, err := routes.Apps.GetFootballSquareGame(getGameParams, resources)

	if err != nil {
		response.ErrorMessage = `Unable to create game`
		responseStr, _ := json.Marshal(response)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(responseStr)

		return
	}

	responseStr, _ := json.Marshal(response)

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseStr)
}

func (routes *HTTPRoutes) getEmptySquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("getEmptySquares Service Acknowledged"))
}

func (routes *HTTPRoutes) reserveSquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)

	var reserveSquareParams app.ReserveSquareParams
	json.NewDecoder(request.Body).Decode(&reserveSquareParams)

	response, err := routes.Apps.ReserveSquare(reserveSquareParams, resources)

	if err != nil {
		response.ErrorMessage = `Unable to reserve square`
		responseStr, _ := json.Marshal(response)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(responseStr)

		return
	}

	responseStr, _ := json.Marshal(response)

	routes.recordSquareReservedTime(reserveSquareParams.GameGUID, resources)

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	writer.WriteHeader(http.StatusOK)
	writer.Write(responseStr)
}

func (routes *HTTPRoutes) saveSquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Save Service Acknowledged"))
}

func (routes *HTTPRoutes) deleteSquare(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("deleteSquare Service Acknowledged"))
}
func (routes *HTTPRoutes) generateNumber(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("GenerateNumber Service Acknowledged"))
}

func (routes *HTTPRoutes) recordSquareReservedTime(gameGUID string, resources *resources.Resources) error {
	redisError := resources.RedisClient.Publish(resources.Context, "SquareReserved:"+gameGUID, "SquareReserved")
	return redisError.Err()
}
