package routes

import (
	"encoding/json"
	"fmt"
	"footballsquaregameservices/app"
	"log"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type RoutesInterface interface {
	Register(resources *resources.Resources) *http.ServeMux
}

type Routes struct {
	Apps app.Game
}

type Handler = func(writer http.ResponseWriter, request *http.Request, resources *resources.Resources)

func NewRoutes() RoutesInterface {
	return &Routes{
		Apps: app.NewGameApp(),
	}
}
func (routes *Routes) Register(resources *resources.Resources) *http.ServeMux {
	log.Println("Registering routes")
	mux := http.NewServeMux()

	routesHandlersMap := map[string]Handler{
		"/":                        routes.home,
		"POST /CreateGame":         routes.createGame,
		"GET /GetGame/{game_guid}": routes.getGame,
		"/GetEmptySquares":         routes.getEmptySquares,
		"/ReserveSquares":          routes.reserveSquares,
		"/SaveSquares":             routes.saveSquares,
		"/DeleteSquare":            routes.deleteSquare,
		"/GenerateNumber":          routes.generateNumber,
	}

	for route, handler := range routesHandlersMap {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, resources)
		})
	}

	return mux
}

func (routes *Routes) home(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	fmt.Fprintf(writer, "{\"Acknowledged\": true}")
}

func (routes *Routes) createGame(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
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

	writer.WriteHeader(http.StatusOK)
	writer.Write(responseStr)
}

func (routes *Routes) getGame(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
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

func (routes *Routes) getEmptySquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("getEmptySquares Service Acknowledged"))
}

func (routes *Routes) reserveSquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("reserveSquares Service Acknowledged"))
}

func (routes *Routes) saveSquares(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Save Service Acknowledged"))
}

func (routes *Routes) deleteSquare(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("deleteSquare Service Acknowledged"))
}
func (routes *Routes) generateNumber(writer http.ResponseWriter, request *http.Request, resources *resources.Resources) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("GenerateNumber Service Acknowledged"))
}
