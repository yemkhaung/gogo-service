package service

import (
	"bytes"
	"fmt"
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

func TestCreateMatch(t *testing.T) {
	client := &http.Client{}
	server := httptest.NewServer(
		http.HandlerFunc(createMatchHandler(formatter)),
	)
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\n  \"players\":  [\n    {\n    \"color\":  \"white\",\n      \"name\": \"bob\"\n    },\n    {\n    \"color\":  \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}
	req.Header.Add("Conetnt-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to creatematchHandler: %v", err)
	}

	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	loc, ok := res.Header["Location"]
	if !ok {
		t.Error("Location header is not set")
	}
	if !strings.Contains(loc[0], "/matches/") {
		t.Error("Location header should contain '/matches/'")
	}
	if len(loc[0]) != len(fakeMatchLocationResult) {
		t.Error("Location value does not contain guid of new match")
	}

	fmt.Printf("Payload: %s", string(payload))
}
