package main

import (
	"context"
	"fmt"
	"footballsquaregameservices/routes"
	"log"
	"net/http"

	"github.com/longvu727/FootballSquaresLibs/DB/db"
	"github.com/longvu727/FootballSquaresLibs/util"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
)

type api struct {
	routes  routes.RoutesInterface
	address string
	server  *http.Server
}

func main() {
	resources, err := getResourcesFromConfigFile("./env", "app", "env")
	if err != nil {
		log.Fatal(err)
	}

	api := &api{
		routes:  routes.NewRoutes(),
		address: fmt.Sprintf(":%s", resources.Config.PORT),
		server:  &http.Server{},
	}

	_ = api.start(resources)
}

func getResourcesFromConfigFile(path string, configName string, configType string) (*resources.Resources, error) {
	config, err := util.LoadConfig(path, configName, configType)
	if err != nil {
		return nil, err
	}

	mysql, err := db.NewMySQL(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resources := resources.NewResources(config, mysql, ctx)

	return resources, nil
}

func (api *api) start(resources *resources.Resources) error {
	log.Printf("Listening on %s", api.address)

	mux := api.routes.Register(resources)

	api.server.Addr = api.address
	api.server.Handler = mux

	err := api.server.ListenAndServe()

	if err != nil {
		log.Print(err)
	}

	return err
}
