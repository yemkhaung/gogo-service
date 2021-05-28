package dronesevents

import (
	"encoding/json"
	"log"
	"runtime/debug"
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
	var telemetryEvent TelemetryEvent
	_ = json.Unmarshal(payload, &telemetryEvent)
	telemetryRecord := convertToTelemetryRecord(telemetryEvent)
	err := tp.QueryStore.UpdateLastTelemetryEvent(telemetryRecord)
	if err != nil {
		log.Printf("Error storing telemetry: %v, err: %s", telemetryRecord, err)
		debug.PrintStack()
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
	var alertEvent AlertEvent
	_ = json.Unmarshal(payload, &alertEvent)
	alertRecord := convertToAlertRecord(alertEvent)
	err := tp.QueryStore.UpdateLastAlertEvent(alertRecord)
	if err != nil {
		log.Printf("Error storing alert event: %v, err: %s", alertRecord, err)
		debug.PrintStack()
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
	var positionEvent PositionEvent
	_ = json.Unmarshal(payload, &positionEvent)
	positionRecord := convertToPositionRecord(positionEvent)
	err := tp.QueryStore.UpdateLastPositionEvent(positionRecord)
	if err != nil {
		log.Printf("Error storing position event: %v, err: %s", positionRecord, err)
		debug.PrintStack()
		return err
	}
	log.Println("Processing position event...[done]")
	return nil
}
