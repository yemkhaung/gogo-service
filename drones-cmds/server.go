package dronescmds

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer returns a Negroni instance
func NewServer() *negroni.Negroni {
	url := parseRabbitURL()
	queue := "telemetry"
	dispatcher := newRabbitMQDispatcher(queue, url)
	return NewServerWithDispatcher(dispatcher)
}

// NewServerWithDispatcher returns Netroni instance using passed dispatcher
func NewServerWithDispatcher(dispatcher amqpDispatcher) *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, dispatcher)
	n.UseHandler(mx)

	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, dispatcher amqpDispatcher) {
	mx.HandleFunc("/api", rootHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/cmds/telemetry", submitTelemetry(formatter, dispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/alerts", submitAlert(formatter, dispatcher)).Methods("POST")
	mx.HandleFunc("/api/cmds/positions", submitPosition(formatter, dispatcher)).Methods("POST")
}

func rootHandler(formatter *render.Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := formatter.JSON(rw, http.StatusOK, struct{ Version string }{"0.0.1"})
		if err != nil {
			log.Printf("Error handling request, %v", err)
		}
	}
}

func parseRabbitURL() string {
	rabbitHost := getenv("RABBITMQ_HOST", "localhost")
	rabbitPort := getenv("RABBITMQ_PORT", "5672")
	return fmt.Sprintf("amqp://guest:guest@%s:%s", rabbitHost, rabbitPort)
}

func getenv(name string, defaultval string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		return defaultval
	}
	return result
}
