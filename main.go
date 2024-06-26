package main

import (
	"fmt"
	"footballsquaregameservices/routes"
	"log"
	"net/http"
	"os"

	"github.com/longvu727/FootballSquaresLibs/util"
)

func main() {
	config, err := util.LoadConfig("./env", "app", "json")
	log.SetOutput(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}

	handler(&config)
}

func handler(config *util.Config) error {

	routes.Register(config)

	address := fmt.Sprintf(":%s", config.PORT)
	log.Printf("Listening on %s", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
