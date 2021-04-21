package main

import (
	"flag"
	"fmt"

	"github.com/codegangsta/negroni"
	catalog "github.com/yemkhaung/gogo-service/backing-catalog"
	fufilment "github.com/yemkhaung/gogo-service/backing-fufilment"
	"github.com/yemkhaung/gogo-service/gogoservice"
)

var (
	server      *negroni.Negroni
	serviceName string
	port        string
)

func init() {
	flag.StringVar(&serviceName, "service", "gogo", "Name of service to run. Must be 'gogo' or 'fufilment' or 'catalog'")
	flag.StringVar(&port, "port", "3000", "Service port.")
}

func main() {
	flag.Parse()
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
