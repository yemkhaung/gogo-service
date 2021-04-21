package catalog

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	serviceURL  = os.Getenv("FUFILMENT_SERVICE_URL")
	servicePort = os.Getenv("FUFILMENT_SERVICE_PORT")
)

// NewServer configures and returns a new server
func NewServer() *negroni.Negroni {
	client := fufilmentWebClient{url: serviceURL + ":" + servicePort}
	return NewServerWithClient(&client)
}

// NewServerWithClient configuresand returns a new server with a client
func NewServerWithClient(client fulfilmentClient) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, client)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, client fulfilmentClient) {
	mx.HandleFunc("/", rootHandler(formatter)).Methods("GET")
	mx.HandleFunc("/catalog", getAllCatalogs(formatter, client)).Methods("GET")
	mx.HandleFunc("/catalog/{sku}", getCataglogDetail(formatter, client)).Methods("GET")
}

func rootHandler(formatter *render.Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := formatter.JSON(rw, http.StatusOK, struct{ Version string }{"0.0.1"})
		if err != nil {
			log.Printf("Error handling request, %v", err)
		}
	}
}
