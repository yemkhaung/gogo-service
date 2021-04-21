package gogoservice

import (
	"testing"

	"github.com/cloudnativego/gogo-engine"
)

func TestAddMatchSucceed(t *testing.T) {
	match := gogo.NewMatch(19, "John", "Bobby")
	repo := newInMemMatchRepository()

	err := repo.addMatch(match)
	if err != nil {
		t.Error("Error adding match to repository. Should add.")
	}

	matches := repo.getMatches()
	if len(matches) != 1 {
		t.Errorf("Expected 1 match in repository. got %d", len(matches))
	}
}

func TestAddMatchEmpty(t *testing.T) {
	repo := newInMemMatchRepository()

	matches := repo.getMatches()
	if len(matches) != 0 {
		t.Errorf("Expected 0 match (empty) in repository. got %d", len(matches))
	}
}
