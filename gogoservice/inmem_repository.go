package gogoservice

import (
	"fmt"

	"github.com/cloudnativego/gogo-engine"
)

type inMemoryMatchRepository struct {
	matches []gogo.Match
}

func newInMemMatchRepository() *inMemoryMatchRepository {
	return &inMemoryMatchRepository{
		matches: []gogo.Match{},
	}
}

func (r *inMemoryMatchRepository) addMatch(m gogo.Match) error {
	r.matches = append(r.matches, m)
	return nil
}

func (r *inMemoryMatchRepository) getMatch(id string) (match gogo.Match, err error) {
	for _, m := range r.matches {
		if m.ID == id {
			return m, nil
		}
	}
	return match, fmt.Errorf("cannot find match: %s", id)
}

func (r *inMemoryMatchRepository) getMatches() ([]gogo.Match, error) {
	return r.matches, nil
}
