package service

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Location", "/matches/"+uuid.NewString())
		err := formatter.JSON(
			rw,
			http.StatusCreated,
			struct{ Test string }{"This is a test"},
		)
		if err != nil {
			log.Printf("Error handling request, %v", err)
		}
	}
}
