package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	catalog "github.com/yemkhaung/gogo-service/backing-catalog"
	fufilment "github.com/yemkhaung/gogo-service/backing-fufilment"
	"github.com/yemkhaung/gogo-service/gogoservice"
)

var (
	server *negroni.Negroni
)

func getenv(name string, defaultval string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		return defaultval
	}
	return result
}

func main() {
	// parse envs
	serviceName := getenv("SERVICE_NAME", "gogo")
	port := getenv("SERVICE_PORT", "3000")

	switch serviceName {
	case "gogo":
		server = gogoservice.NewServer()

	case "fufilment":
		server = fufilment.NewServer()

	case "catalog":
		server = catalog.NewServer()
	}
	fmt.Printf("Running '%s' service...\n", serviceName)
	server.Run(":" + port)
}
