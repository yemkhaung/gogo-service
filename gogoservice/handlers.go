package gogoservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloudnativego/gogo-engine"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		payload, _ := ioutil.ReadAll(r.Body)
		var newMatchRquest newMatchRequest
		err := json.Unmarshal(payload, &newMatchRquest)
		if err != nil {
			log.Printf("Error parsing request payload, %v\n", err)
			_ = formatter.Text(rw, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if !newMatchRquest.isValid() {
			log.Printf("Invalid match request body: %v\n", newMatchRquest)
			_ = formatter.Text(rw, http.StatusBadRequest, "Invalid request")
			return
		}

		newMatch := gogo.NewMatch(newMatchRquest.GridSize, newMatchRquest.PlayerBlack, newMatchRquest.PlayerWhite)
		err = repo.addMatch(newMatch)
		if err != nil {
			log.Println(err)
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error adding match to repository")
			return
		}

		rw.Header().Add("Location", "/matches/"+newMatch.ID)
		_ = formatter.JSON(
			rw,
			http.StatusCreated,
			newMatchResponse{
				ID:          newMatch.ID,
				GridSize:    newMatch.GridSize,
				PlayerBlack: newMatch.PlayerBlack,
				PlayerWhite: newMatch.PlayerWhite,
			},
		)
	}
}

func getOneMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		matchID := vars["id"]
		match, err := repo.getMatch(matchID)
		if err != nil {
			log.Println(err)
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error retrieving match from repository")
			return
		}
		_ = formatter.JSON(
			rw,
			http.StatusOK,
			match,
		)
	}
}

func getMatchesHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		matches, err := repo.getMatches()
		if err != nil {
			log.Println(err)
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error retrieving matches from repository")
			return
		}
		_ = formatter.JSON(
			rw,
			http.StatusOK,
			matches,
		)
	}
}
