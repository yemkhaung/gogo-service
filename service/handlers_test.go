package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

//nolint // not required for test case
func TestCreateMatch(t *testing.T) {
	// prepare
	client := &http.Client{}
	repo := newInMemMatchRepository()
	server := httptest.NewServer(
		http.HandlerFunc(createMatchHandler(formatter, repo)),
	)
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\n  \"players\":  [\n    {\n    \"color\":  \"white\",\n      \"name\": \"bob\"\n    },\n    {\n    \"color\":  \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}
	req.Header.Add("Conetnt-Type", "application/json")

	// execute
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to creatematchHandler: %v", err)
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	// assert Status code
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	// assert Location header
	loc, ok := res.Header["Location"]
	if !ok || loc == nil {
		t.Error("Location header is not set")
	}

	// assert Location header patttern
	if !strings.Contains(loc[0], "/matches/") {
		t.Error("Location header should contain '/matches/'")
	}

	// assert Location header length
	if len(loc[0]) != len(fakeMatchLocationResult) {
		t.Errorf("Expected Location header to have same length of %d, got %d", len(fakeMatchLocationResult), len(loc[0]))
	}

	// assert ID in response with ID in Location Header
	var matchResponse newMatchResponse
	err = json.Unmarshal(payload, &matchResponse)
	if err != nil {
		t.Error("Error. Response payload cannot be JSON encoded.")
	}
	if matchResponse.ID == "" || !strings.Contains(loc[0], matchResponse.ID) {
		t.Errorf("Expected Match ID of %s, got %s", loc[0], matchResponse.ID)
	}

	// assert match is added to repository
	matches := repo.getMatches()
	if len(matches) != 1 {
		t.Errorf("Expected 1 match exists in repository, got %d", len(matches))
	}
	match := matches[0]
	if match.GridSize != matchResponse.GridSize {
		t.Errorf("Expected matching gridsize of %d, got %d", match.GridSize, matchResponse.GridSize)
	}
	if len(matchResponse.Players) != 2 {
		t.Error("Expected 2 players from added match.")
	}
	for _, player := range matchResponse.Players {
		if player.Color == "black" {
			if player.Name != match.PlayerBlack {
				t.Errorf("Expected black player name of %s, got %s", player.Name, match.PlayerBlack)
			}
		} else if player.Color == "white" {
			if player.Name != match.PlayerWhite {
				t.Errorf("Expected white player name of %s, got %s", player.Name, match.PlayerWhite)
			}
		}
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
	client := &http.Client{}
	repo := newInMemMatchRepository()
	server := httptest.NewServer(
		http.HandlerFunc(createMatchHandler(formatter, repo)),
	)
	defer server.Close()

	for _, val := range invalidRequestTable {
		body := []byte(val)
		req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
		}
		req.Header.Add("Conetnt-Type", "application/json")
		// execute
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Error in POST to creatematchHandler: %v", err)
		}
		defer res.Body.Close()
		// assert Status code
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected response status 400, received %s", res.Status)
		}
	}
}
