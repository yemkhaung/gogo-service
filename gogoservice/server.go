package gogoservice

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a new server
func NewServer() *negroni.Negroni {
	dbURL := os.Getenv("MONGODB_URL")
	return NewServerWithRepo(newPersistRepository(dbURL))
}

// NewServerWithRepo configures a new server with a repository instance
func NewServerWithRepo(repo matchRepository) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, repo)
	n.UseHandler(mx)

	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo matchRepository) {
	mx.HandleFunc("/", rootHandler(formatter)).Methods("GET")
	mx.HandleFunc("/matches", createMatchHandler(formatter, repo)).Methods("POST")
	mx.HandleFunc("/matches", getMatchesHandler(formatter, repo)).Methods("GET")
	mx.HandleFunc("/matches/{id}", getOneMatchHandler(formatter, repo)).Methods("GET")
}

func rootHandler(formatter *render.Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := formatter.JSON(rw, http.StatusOK, struct{ Version string }{"0.0.1"})
		if err != nil {
			log.Printf("Error handling request, %v", err)
		}
	}
}
