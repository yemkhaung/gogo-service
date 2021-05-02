package dronesevents

import (
	"encoding/json"
	"log"
)

type eventProcessor interface {
	Process(event interface{}) error
}

type telemetryProcessor struct {
	QueryStore eventRepository
}

func (tp *telemetryProcessor) Process(event interface{}) error {
	// TODO do complex calculations
	log.Println("Processing telemetry event...[started]")
	payload, _ := json.Marshal(event)
	var telemetryEvent TelemetryUpdatedEvent
	_ = json.Unmarshal(payload, &telemetryEvent)

	err := tp.QueryStore.UpdateLastTelemetryEvent(telemetryEvent)
	if err != nil {
		log.Printf("Error storing telemetry event: %s", err)
		return err
	}
	log.Println("Processing telemetry event...[done]")
	return nil
}

type alertsProcessor struct {
	QueryStore eventRepository
}

func (tp *alertsProcessor) Process(event interface{}) error {
	// TODO do complex calculations
	log.Println("Processing alert event...[started]")
	payload, _ := json.Marshal(event)
	var alertEvent AlertSignalledEvent
	_ = json.Unmarshal(payload, &alertEvent)

	err := tp.QueryStore.UpdateLastAlertEvent(alertEvent)
	if err != nil {
		log.Printf("Error storing alert event: %s", err)
		return err
	}
	log.Println("Processing alert event...[done]")
	return nil
}

type positionsProcessor struct {
	QueryStore eventRepository
}

func (tp *positionsProcessor) Process(event interface{}) error {
	// TODO do complex calculations
	log.Println("Processing position event...[started]")
	payload, _ := json.Marshal(event)
	var positionEvent PositionChangedEvent
	_ = json.Unmarshal(payload, &positionEvent)

	err := tp.QueryStore.UpdateLastPositionEvent(positionEvent)
	if err != nil {
		log.Printf("Error storing position event: %s", err)
		return err
	}
	log.Println("Processing position event...[done]")
	return nil
}
