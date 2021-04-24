package gogoservice

import "github.com/cloudnativego/gogo-engine"

type newMatchRequest struct {
	GridSize    int    `json:"gridsize"`
	PlayerBlack string `json:"player_black"`
	PlayerWhite string `json:"player_white"`
}

type newMatchResponse struct {
	ID          string `json:"id"`
	GridSize    int    `json:"gridsize"`
	StartedAt   string `json:"started_at"`
	PlayerBlack string `json:"player_black"`
	PlayerWhite string `json:"player_white"`
}

type matchRepository interface {
	addMatch(m gogo.Match) error
	getMatches() ([]gogo.Match, error)
	getMatch(id string) (gogo.Match, error)
}

func (r *newMatchRequest) isValid() bool {
	// validate grid size
	if r.GridSize != 19 && r.GridSize != 13 && r.GridSize != 9 {
		return false
	}
	// validate players info
	if r.PlayerBlack == "" || r.PlayerWhite == "" {
		return false
	}
	return true
}
