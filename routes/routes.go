package routes

import (
	"encoding/json"
	"fmt"
	"footballsquaregameservices/app"
	"log"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/util"
)

type Handler = func(writer http.ResponseWriter, request *http.Request, config *util.Config)

func Register(config *util.Config) {
	log.Println("Registering routes")
	routes := map[string]Handler{
		"/":                        home,
		"POST /CreateGame":         createGame,
		"GET /GetGame/{game_guid}": getGame,
		"/GetEmptySquares":         getEmptySquares,
		"/ReserveSquares":          reserveSquares,
		"/SaveSquares":             saveSquares,
		"/DeleteSquare":            deleteSquare,
		"/GenerateNumber":          generateNumber,
	}

	for route, handler := range routes {
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, config)
		})
	}
}

func home(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	fmt.Fprintf(writer, "{\"Acknowledged\": true}")
}

func createGame(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)

	response, err := app.CreateFootballSquareGame(request, config)

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

func getGame(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	response, err := app.GetFootballSquareGame(request, config)

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

func getEmptySquares(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("getEmptySquares Service Acknowledged"))
}

func reserveSquares(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("reserveSquares Service Acknowledged"))
}

func saveSquares(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Save Service Acknowledged"))
}

func deleteSquare(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("deleteSquare Service Acknowledged"))
}
func generateNumber(writer http.ResponseWriter, request *http.Request, config *util.Config) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("GenerateNumber Service Acknowledged"))
}
