package routes

import (
	"encoding/json"
	"fmt"
	"footballsquaregameservices/app"
	"log"
	"net/http"
)

type Handler = func(writer http.ResponseWriter, request *http.Request)

func Register() {
	log.Println("Registering routes")
	routes := map[string]Handler{
		"/":                home,
		"/CreateGame":      createGame,
		"/GetGame":         getGame,
		"/GetEmptySquares": getEmptySquares,
		"/ReserveSquares":  reserveSquares,
		"/SaveSquares":     saveSquares,
		"/DeleteSquare":    deleteSquare,
		"/GenerateNumber":  generateNumber,
	}

	for route, handler := range routes {
		http.HandleFunc(route, handler)
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "{\"Acknowledged\": true}")
}

func createGame(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)

	response, err := app.CreateFootballSquareGame(request)

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

func getGame(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("getGame Service Acknowledged"))
}

func getEmptySquares(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("getEmptySquares Service Acknowledged"))
}

func reserveSquares(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("reserveSquares Service Acknowledged"))
}

func saveSquares(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Save Service Acknowledged"))
}

func deleteSquare(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("deleteSquare Service Acknowledged"))
}
func generateNumber(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request for %s\n", request.URL.Path)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("GenerateNumber Service Acknowledged"))
}
