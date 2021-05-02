package dronesevents

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures a new server with event-consumers and event-processors
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	server := negroni.Classic()

	m := mux.NewRouter()
	initRoutes(formatter, m)
	server.UseHandler(m)

	mongoURL := parseMongoURL()
	registry := newConsumerRegistry()
	repository := &mongoEventRepository{URL: mongoURL}

	initRegistry(registry, repository)
	err := registry.consumeEvents()
	if err != nil {
		log.Fatalf("Error consuming events: %s", err)
	}

	return server
}

func initRoutes(formatter *render.Render, m *mux.Router) {
	m.HandleFunc("/", rootHandler(formatter)).Methods("GET")
}

func initRegistry(registry *consumerRegistry, repository *mongoEventRepository) {
	rabbitURL := parseRabbitURL()
	// Consumers
	registry.RegisterConsumer("telemetry", &amqpConsumer{AmqpURL: rabbitURL, QueueName: "telemetry"})
	registry.RegisterConsumer("alert", &amqpConsumer{AmqpURL: rabbitURL, QueueName: "alert"})
	registry.RegisterConsumer("position", &amqpConsumer{AmqpURL: rabbitURL, QueueName: "position"})
	// Processors
	registry.RegisterProcessor("telemetry", &telemetryProcessor{QueryStore: repository})
	registry.RegisterProcessor("alert", &alertsProcessor{QueryStore: repository})
	registry.RegisterProcessor("position", &positionsProcessor{QueryStore: repository})
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

func parseMongoURL() string {
	mongoHost := getenv("MONGODB_HOST", "localhost")
	mongoPort := getenv("MONGODB_PORT", "27017")
	return fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
}

func getenv(name string, defaultval string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		return defaultval
	}
	return result
}
