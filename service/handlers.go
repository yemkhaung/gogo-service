package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloudnativego/gogo-engine"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		payload, _ := ioutil.ReadAll(r.Body)
		var newMatchRquest newMatchRequest
		err := json.Unmarshal(payload, &newMatchRquest)
		if err != nil {
			_ = formatter.Text(rw, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if !newMatchRquest.isValid() {
			_ = formatter.Text(rw, http.StatusBadRequest, "Invalid request")
			return
		}

		playerBlack, playerWhite := parsePlayerNames(newMatchRquest.Players)
		newMatch := gogo.NewMatch(newMatchRquest.GridSize, playerBlack, playerWhite)
		err = repo.addMatch(newMatch)
		if err != nil {
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error adding match to repository")
			return
		}

		rw.Header().Add("Location", "/matches/"+newMatch.ID)
		_ = formatter.JSON(
			rw,
			http.StatusCreated,
			newMatchResponse{
				ID:       newMatch.ID,
				GridSize: newMatch.GridSize,
				Players:  createPlayers(playerBlack, playerWhite),
			},
		)
	}
}
