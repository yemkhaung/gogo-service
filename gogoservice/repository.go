package gogoservice

import "github.com/cloudnativego/gogo-engine"

type inMemoryMatchRepository struct {
	matches []gogo.Match
}

func newInMemMatchRepository() *inMemoryMatchRepository {
	return &inMemoryMatchRepository{
		matches: []gogo.Match{},
	}
}

func (r *inMemoryMatchRepository) addMatch(m gogo.Match) (err error) {
	r.matches = append(r.matches, m)
	return err
}

func (r *inMemoryMatchRepository) getMatches() []gogo.Match {
	return r.matches
}
