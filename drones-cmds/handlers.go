package dronescmds

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/unrolled/render"
)

type telemetryCommandEvent struct {
	Telemetry  telemetryCommand `json:"telemetry"`
	ReceivedOn int64            `json:"received_on"`
}

type alertCommandEvent struct {
	Alert      alertCommand `json:"alert"`
	ReceivedOn int64        `json:"received_on"`
}

type positionCommandEvent struct {
	Position   positionCommand `json:"position"`
	ReceivedOn int64           `json:"received_on"`
}

func submitTelemetry(formatter *render.Render, dispatcher amqpDispatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// parse message
		paylod, _ := ioutil.ReadAll(req.Body)
		var cmd telemetryCommand
		err := json.Unmarshal(paylod, &cmd)
		if err != nil || !cmd.isValid() {
			_ = formatter.Text(rw, http.StatusBadRequest, "Invalid request content")
		}

		// dispatch message to event-store
		queue := "telemetry"
		event := telemetryCommandEvent{
			Telemetry:  cmd,
			ReceivedOn: time.Now().Unix(),
		}
		err = dispatcher.Dispatch(queue, event)
		if err != nil {
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error dispatching event to queue")
		}

		_ = formatter.JSON(rw, http.StatusCreated, event)
	}
}

func submitAlert(formatter *render.Render, dispatcher amqpDispatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// parse message
		paylod, _ := ioutil.ReadAll(req.Body)
		var cmd alertCommand
		err := json.Unmarshal(paylod, &cmd)
		if err != nil || !cmd.isValid() {
			_ = formatter.Text(rw, http.StatusBadRequest, "Invalid request content")
		}

		// dispatch message to event-store
		queue := "alert"
		event := alertCommandEvent{
			Alert:      cmd,
			ReceivedOn: time.Now().Unix(),
		}
		err = dispatcher.Dispatch(queue, event)
		if err != nil {
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error dispatching event to queue")
		}

		_ = formatter.JSON(rw, http.StatusCreated, event)
	}
}

func submitPosition(formatter *render.Render, dispatcher amqpDispatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// parse message
		paylod, _ := ioutil.ReadAll(req.Body)
		var cmd positionCommand
		err := json.Unmarshal(paylod, &cmd)
		if err != nil || !cmd.isValid() {
			_ = formatter.Text(rw, http.StatusBadRequest, "Invalid request content")
		}

		// dispatch message to event-store
		queue := "position"
		event := positionCommandEvent{
			Position:   cmd,
			ReceivedOn: time.Now().Unix(),
		}
		err = dispatcher.Dispatch(queue, event)
		if err != nil {
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error dispatching event to queue")
		}

		_ = formatter.JSON(rw, http.StatusCreated, event)
	}
}
