package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"squareservices/routes"
	"squareservices/util"
)

func main() {
	config, err := util.LoadConfig("./env", "app", "env")
	log.SetOutput(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}

	handler(config)
}

func handler(config util.Config) error {

	routes.Register()

	address := fmt.Sprintf(":%s", config.PORT)
	log.Printf("Listening on %s", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
