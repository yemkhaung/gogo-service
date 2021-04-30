package dronescmds

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeDispatcher struct {
	queue []interface{}
}

func (q *fakeDispatcher) Dispatch(queue string, event interface{}) (err error) {
	q.queue = append(q.queue, event)
	return
}

func TestSubmitTelemetrySucceed(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\": \"d-001\", \"uptime\": 100}")
	req, err := http.NewRequest("POST", "/api/cmds/telemetry", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code of %d, got %d", http.StatusCreated, recorder.Code)
		return
	}
	payload, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
		return
	}
	var responseEvent telemetryCommandEvent
	err = json.Unmarshal(payload, &responseEvent)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err)
		return
	}
	if responseEvent.Telemetry.DroneID != "d-001" || responseEvent.Telemetry.Uptime != 100 || responseEvent.ReceivedOn == 0 {
		t.Errorf("Expected drone_id: d-001 got %s and uptime:100 got %d", responseEvent.Telemetry.DroneID, responseEvent.Telemetry.Uptime)
		return
	}
	if len(dispatcher.queue) != 1 {
		t.Errorf("Expected queue size of 1, got %d", len(dispatcher.queue))
		return
	}
}

func TestSubmitTelemetryIvalidFormat(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\": \"\", \"uptime\": 0}")
	req, err := http.NewRequest("POST", "/api/cmds/telemetry", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code of %d, got %d", http.StatusBadRequest, recorder.Code)
		return
	}
}

func TestSubmitAlertSucceed(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\": \"d-001\", \"fault_code\": 500}")
	req, err := http.NewRequest("POST", "/api/cmds/alerts", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code of %d, got %d", http.StatusCreated, recorder.Code)
		return
	}
	payload, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
		return
	}
	var responseEvent alertCommandEvent
	err = json.Unmarshal(payload, &responseEvent)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err)
		return
	}
	if responseEvent.Alert.DroneID != "d-001" || responseEvent.Alert.FaultCode != 500 || responseEvent.ReceivedOn == 0 {
		t.Errorf("Expected drone_id:'d-001' got %s and fault_code:500 got %d", responseEvent.Alert.DroneID, responseEvent.Alert.FaultCode)
		return
	}
	if len(dispatcher.queue) != 1 {
		t.Errorf("Expected queue size of 1, got %d", len(dispatcher.queue))
		return
	}
}

func TestSubmitAlertIvalidFormat(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\": \"\", \"fault_code\": 0}")
	req, err := http.NewRequest("POST", "/api/cmds/alerts", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code of %d, got %d", http.StatusBadRequest, recorder.Code)
		return
	}
}

func TestSubmitPositionSucceed(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\":\"d-001\", \"latitude\": 81.231, \"longitude\": 43.1231, \"altitude\": 2301.1, \"current_speed\": 41.3, \"heading_cardinal\": 1}")
	req, err := http.NewRequest("POST", "/api/cmds/positions", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code of %d, got %d", http.StatusCreated, recorder.Code)
		return
	}
	payload, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
		return
	}
	var responseEvent positionCommandEvent
	err = json.Unmarshal(payload, &responseEvent)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err)
		return
	}
	if responseEvent.Position.DroneID != "d-001" || responseEvent.Position.Altitude != 2301.1 || responseEvent.ReceivedOn == 0 {
		t.Errorf("Expected drone_id:'d-001' got %s and altitude:2301.1 got %f", responseEvent.Position.DroneID, responseEvent.Position.Altitude)
		return
	}
	if len(dispatcher.queue) != 1 {
		t.Errorf("Expected queue size of 1, got %d", len(dispatcher.queue))
		return
	}
}

func TestSubmitPositionIvalidFormat(t *testing.T) {
	dispatcher := &fakeDispatcher{
		queue: []interface{}{},
	}
	server := NewServerWithDispatcher(dispatcher)

	recorder := httptest.NewRecorder()
	body := []byte("{\"drone_id\": \"\"}")
	req, err := http.NewRequest("POST", "/api/cmds/alerts", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating new request: %s", err)
		return
	}
	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code of %d, got %d", http.StatusBadRequest, recorder.Code)
		return
	}
}
