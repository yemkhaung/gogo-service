package gogoservice

import (
	"testing"

	"github.com/cloudnativego/gogo-engine"
)

func beforeEach() []matchRepository {
	repos := []matchRepository{
		newInMemMatchRepository(),
		// newPersistRepository(os.Getenv("MONGODB_URL")), # for integration test
	}
	return repos
}

func TestAddMatchSucceed(t *testing.T) {
	repos := beforeEach()
	for _, repo := range repos {
		match := gogo.NewMatch(19, "John", "Bobby")
		err := repo.addMatch(match)
		if err != nil {
			t.Errorf("Error adding match to repository, %v", err)
		}
		matches, err := repo.getMatches()
		if err != nil || len(matches) != 1 {
			t.Errorf("Expected 1 match in repository. got %d", len(matches))
		}
	}
}

func TestAddMatchEmpty(t *testing.T) {
	repos := beforeEach()
	for _, repo := range repos {
		matches, err := repo.getMatches()
		if err != nil || len(matches) != 0 {
			t.Errorf("Expected 0 match (empty) in repository. got %d", len(matches))
		}
	}
}
