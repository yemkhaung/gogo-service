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

// Getenv gets environment variable with default value
func Getenv(name string, defaultval string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		return defaultval
	}
	return result
}

func main() {
	// parse envs
	serviceName := Getenv("SERVICE_NAME", "gogo")
	port := Getenv("SERVICE_PORT", "3000")

	switch serviceName {
	case "gogo":
		dbURL := Getenv("MONGODB_URL", "mongodb://localhost:27017")
		server = gogoservice.NewServer(dbURL)

	case "fufilment":
		server = fufilment.NewServer()

	case "catalog":
		server = catalog.NewServer()
	}
	fmt.Printf("Running '%s' service...\n", serviceName)
	server.Run(":" + port)
}
