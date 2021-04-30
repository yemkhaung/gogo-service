package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	catalog "github.com/yemkhaung/gogo-service/backing-catalog"
	fufilment "github.com/yemkhaung/gogo-service/backing-fufilment"
	dronescmds "github.com/yemkhaung/gogo-service/drones-cmds"
	gogo "github.com/yemkhaung/gogo-service/gogoservice"
)

func getenv(name string, defaultval string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		return defaultval
	}
	return result
}

func getServiceEnvs() (serviceName string, servicePort string) {
	serviceName = getenv("SERVICE_NAME", "gogo")
	servicePort = getenv("SERVICE_PORT", "3000")
	return serviceName, servicePort
}

func main() {
	serviceName, servicePort := getServiceEnvs()
	var server *negroni.Negroni

	switch serviceName {
	case "gogo":
		server = gogo.NewServer()

	case "fufilment":
		server = fufilment.NewServer()

	case "catalog":
		server = catalog.NewServer()

	case "dronescmds":
		server = dronescmds.NewServer()
	}

	fmt.Printf("Running '%s' service...\n", serviceName)
	server.Run(":" + servicePort)
}
