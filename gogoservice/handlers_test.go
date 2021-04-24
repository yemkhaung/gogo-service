package gogoservice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

//nolint // not required for test case
func TestCreateMatch(t *testing.T) {
	// prepare
	repo := newInMemMatchRepository()
	server := NewServerWithRepo(repo)

	body := []byte("{\n  \"gridsize\": 19,\n  \"player_white\": \"bob\",\n  \"player_black\": \"alfred\"\n}")

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/matches", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	server.ServeHTTP(resp, req)

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
		return
	}

	// assert Status code
	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response status 201, received %d", resp.Code)
		return
	}

	// assert Location header
	loc := resp.Header().Get("Location")
	if loc == "" {
		t.Error("Location header is not set")
		return
	}

	// assert Location header patttern
	if !strings.Contains(loc, "/matches/") {
		t.Error("Location header should contain '/matches/'")
		return
	}

	// assert Location header length
	if len(loc) != len(fakeMatchLocationResult) {
		t.Errorf("Expected Location header to have same length of %d, got %d", len(fakeMatchLocationResult), len(loc))
		return
	}

	// assert ID in response with ID in Location Header
	var matchResponse newMatchResponse
	err = json.Unmarshal(payload, &matchResponse)
	if err != nil {
		t.Error("Error. Response payload cannot be JSON encoded.")
		return
	}
	if matchResponse.ID == "" || !strings.Contains(loc, matchResponse.ID) {
		t.Errorf("Expected Match ID of %s, got %s", loc, matchResponse.ID)
		return
	}

	// assert match is added to repository
	matches, err := repo.getMatches()
	if err != nil || len(matches) != 1 {
		t.Errorf("Expected 1 match exists in repository, got %d", len(matches))
		return
	}
	match := matches[0]
	if match.ID != matchResponse.ID {
		t.Errorf("Expected match ID of %s, got %s", match.ID, matchResponse.ID)
		return
	}
	if match.GridSize != matchResponse.GridSize {
		t.Errorf("Expected matching gridsize of %d, got %d", match.GridSize, matchResponse.GridSize)
		return
	}
	if match.PlayerBlack != matchResponse.PlayerBlack {
		t.Errorf("Expected black player name of %s, got %s", matchResponse.PlayerBlack, match.PlayerBlack)
		return
	}
	if match.PlayerWhite != matchResponse.PlayerWhite {
		t.Errorf("Expected white player name of %s, got %s", matchResponse.PlayerWhite, match.PlayerWhite)
		return
	}
}

var invalidRequestTable = []string{
	"{\n  \"gridsize\": 200,\n  \"players\":  [\n    {\n    \"color\":  \"white\",\n      \"name\": \"bob\"\n    },\n    {\n    \"color\":  \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}",
	"{\n  \"gridsize\": 19,\n  \"players\":  []\n}",
	"{\n  \"gridsize\": 19,\n  \"players\":  [\n    {\n    \"color\":  \"yellow\",\n      \"name\": \"bob\"\n    },\n    {\n    \"color\":  \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}",
	"{\n  \"gridsize\": 19,\n  \"players\":  [\n    {\n    \"color\":  \"white\",\n      \"name\": \"bob\"\n    }\n  ]\n}",
}

func TestNewMatchInvalidRequest(t *testing.T) {
	// prepare
	repo := newInMemMatchRepository()
	server := NewServerWithRepo(repo)

	for _, val := range invalidRequestTable {
		body := []byte(val)
		req, err := http.NewRequest("POST", "/matches", bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
		}
		req.Header.Add("Conetnt-Type", "application/json")
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		// assert Status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("Expected response status 400, received %d", res.Code)
		}
	}
}
