package service

import "github.com/cloudnativego/gogo-engine"

type newMatchRequest struct {
	GridSize int      `json:"gridsize"`
	Players  []player `json:"players"`
}

type newMatchResponse struct {
	ID        string   `json:"id"`
	GridSize  int      `json:"gridsize"`
	StartedAt string   `json:"started_at"`
	Players   []player `json:"players"`
}

type player struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type matchRepository interface {
	addMatch(m gogo.Match) (err error)
	getMatches() []gogo.Match
}

func (r *newMatchRequest) isValid() bool {
	// validate grid size
	if r.GridSize != 19 && r.GridSize != 13 && r.GridSize != 9 {
		return false
	}
	// validate players info
	if len(r.Players) != 2 {
		return false
	}
	for _, p := range r.Players {
		if p.Color != "black" && p.Color != "white" {
			return false
		}
	}
	return true
}

func parsePlayerNames(plist []player) (playerBlack string, playerWhite string) {
	for _, player := range plist {
		if player.Color == "black" {
			playerBlack = player.Name
		} else if player.Color == "white" {
			playerWhite = player.Name
		}
	}
	return playerBlack, playerWhite
}

func createPlayers(playerBlack string, playerWhite string) []player {
	return []player{
		{Color: "black", Name: playerBlack},
		{Color: "white", Name: playerWhite},
	}
}
