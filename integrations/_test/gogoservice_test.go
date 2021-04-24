package integrations

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cloudnativego/gogo-engine"
)

type newMatchResponse struct {
	ID          string `json:"id"`
	GridSize    int    `json:"gridsize"`
	StartedAt   string `json:"started_at"`
	PlayerBlack string `json:"player_black"`
	PlayerWhite string `json:"player_white"`
}

var (
	baseURL = "http://localhost:3000/matches"
	matches = [][]byte{
		[]byte("{\n  \"gridsize\": 19,\n  \"player_white\": \"bob\",\n  \"player_black\": \"alfred\"\n}"),
		[]byte("{\n  \"gridsize\": 19,\n  \"player_white\": \"john\",\n  \"player_black\": \"mary\"\n}"),
	}
)

func TestIntegration(t *testing.T) {
	emptyMatches, err := getMatches(t)
	if err != nil {
		t.Error("Error getting matches")
		return
	}
	if len(emptyMatches) != 0 {
		t.Errorf("Expected empty matches, got %d", len(emptyMatches))
		return
	}

	// assert 1st match
	resp1, err := addMatch(t, matches[0])
	if err != nil {
		t.Error("Error adding matches")
		return
	}
	if resp1.PlayerBlack != "alfred" {
		t.Errorf("Expected playerBlack: 'alfred', got %s", resp1.PlayerBlack)
		return
	}
	matchList, err := getMatches(t)
	if err != nil {
		t.Error("Error getting matches")
		return
	}
	if len(matchList) != 1 {
		t.Errorf("Expected one match added, got %d", len(matchList))
		return
	}
	if matchList[0].PlayerWhite != "bob" {
		t.Errorf("Expected playerWhite: 'bob', got %s", matchList[0].PlayerWhite)
		return
	}

	// assert 2nd match
	resp2, err := addMatch(t, matches[1])
	if err != nil {
		t.Error("Error adding matches")
		return
	}
	if resp2.PlayerBlack != "mary" {
		t.Errorf("Expected playerBlack: 'mary', got %s", resp2.PlayerBlack)
		return
	}
	matchList, err = getMatches(t)
	if err != nil {
		t.Error("Error getting matches")
		return
	}
	if len(matchList) != 2 {
		t.Errorf("Expected one match added, got %d", len(matchList))
		return
	}
	if matchList[1].PlayerWhite != "john" {
		t.Errorf("Expected playerWhite: 'john', got %s", matchList[0].PlayerWhite)
		return
	}

	// assert match details
	mDetails1, err := getMatchDetails(t, resp1.ID)
	if err != nil {
		t.Error("Error getting match details")
		return
	}
	if mDetails1.GridSize != resp1.GridSize {
		t.Errorf("Expected grid size of %d, got %d", resp1.GridSize, mDetails1.GridSize)
		return
	}
	mDetails2, err := getMatchDetails(t, resp2.ID)
	if err != nil {
		t.Error("Error getting match details")
		return
	}
	if mDetails2.PlayerWhite != resp2.PlayerWhite {
		t.Errorf("Expected player-white name of %s, got %s", resp2.PlayerWhite, mDetails2.PlayerWhite)
		return
	}

}

// helper function for adding match to repository
func addMatch(t *testing.T, matchBody []byte) (matchResp newMatchResponse, err error) {

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(matchBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Errorf("Error when adding match. Response: %v", resp)
		return matchResp, err
	}
	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	_ = json.Unmarshal(payload, &matchResp)
	return matchResp, nil
}

// helper function for getting matches from repository
func getMatches(t *testing.T) (matches []gogo.Match, err error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Error("Error when sending request to the server")
		return matches, err
	}
	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	_ = json.Unmarshal(payload, &matches)
	return matches, nil
}

// helper function for getting match details from repository
func getMatchDetails(t *testing.T, id string) (match gogo.Match, err error) {
	resp, err := http.Get(baseURL + "/" + id)
	if err != nil {
		t.Error("Errored when sending request to the server")
		return match, err
	}
	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	_ = json.Unmarshal(payload, &match)
	return match, nil
}
