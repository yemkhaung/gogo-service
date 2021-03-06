package dronesevents

import (
	"testing"
	"time"
)

type fakeRepository struct {
	lastTelemetry TelemetryRecord
	lastAlert     AlertRecord
	lastPosition  PositionRecord
}

func (repo *fakeRepository) UpdateLastTelemetryEvent(record TelemetryRecord) error {
	repo.lastTelemetry = record
	return nil
}

func (repo *fakeRepository) UpdateLastAlertEvent(record AlertRecord) error {
	repo.lastAlert = record
	return nil
}
func (repo *fakeRepository) UpdateLastPositionEvent(record PositionRecord) error {
	repo.lastPosition = record
	return nil
}

type fakeConsumer struct {
	channel chan interface{}
}

func (fc *fakeConsumer) Dequeue() (<-chan interface{}, error) {
	return fc.channel, nil
}

func TestConsumeTelemetryEvents(t *testing.T) {
	// declare fake stubs
	queryStore := &fakeRepository{}
	telemetryConsumer := &fakeConsumer{channel: make(chan interface{})}
	alertsConsumer := &fakeConsumer{channel: make(chan interface{})}
	positionsConsumer := &fakeConsumer{channel: make(chan interface{})}
	// register consumers
	registry := newConsumerRegistry()
	registry.RegisterConsumer("telemetry", telemetryConsumer)
	registry.RegisterConsumer("alert", alertsConsumer)
	registry.RegisterConsumer("position", positionsConsumer)
	registry.RegisterProcessor("telemetry", &telemetryProcessor{QueryStore: queryStore})
	registry.RegisterProcessor("alert", &alertsProcessor{QueryStore: queryStore})
	registry.RegisterProcessor("position", &positionsProcessor{QueryStore: queryStore})
	// start consumers
	err := registry.consumeEvents()
	if err != nil {
		t.Errorf("Error consuming events: %s", err)
	}
	// testing events
	telemetryEvent := TelemetryEvent{
		Telemetry: Telemetry{
			DroneID:          "d-001",
			RemainingBattery: 50,
			Uptime:           120,
			CoreTemp:         25,
		},
		ReceivedOn: time.Now().Unix(),
	}
	alertEvent := AlertEvent{
		Alert: Alert{
			DroneID:   "d-001",
			FaultCode: 55,
		},
		ReceivedOn: time.Now().Unix(),
	}
	alertEvent2 := AlertEvent{
		Alert: Alert{
			DroneID:   "d-001",
			FaultCode: 105,
		},
		ReceivedOn: time.Now().Unix(),
	}
	positionEvent := PositionEvent{
		Position: Position{
			DroneID:         "d-001",
			Latitude:        100.5,
			Longitude:       50.5,
			Altitude:        200.5,
			CurrentSpeed:    20.5,
			HeadingCardinal: 10,
		},
		ReceivedOn: time.Now().Unix(),
	}
	// publish test events
	telemetryConsumer.channel <- telemetryEvent
	alertsConsumer.channel <- alertEvent
	alertsConsumer.channel <- alertEvent2
	positionsConsumer.channel <- positionEvent

	// wait for events to be processed and be stored in query-store
	time.Sleep(time.Second * 1)

	// assertions
	if queryStore.lastTelemetry.DroneID != telemetryEvent.Telemetry.DroneID {
		t.Errorf("Expected telemetry DroneID: %s, got %s", telemetryEvent.Telemetry.DroneID, queryStore.lastTelemetry.DroneID)
	}
	if queryStore.lastAlert.FaultCode != alertEvent2.Alert.FaultCode {
		t.Errorf("Expected alert FaultCode: %d, got %d", alertEvent2.Alert.FaultCode, queryStore.lastAlert.FaultCode)
	}
	if queryStore.lastPosition.Altitude != positionEvent.Position.Altitude {
		t.Errorf("Expected position Altitude: %f, got %f", positionEvent.Position.Altitude, queryStore.lastPosition.Altitude)
	}
}
