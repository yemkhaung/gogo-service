package dronesevents

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TelemetryRecord for receiving drone telemetry data
type TelemetryRecord struct {
	RecordID         primitive.ObjectID `bson:"_id,omitempty"`
	DroneID          string             `bson:"drone_id"`
	RemainingBattery int                `bson:"battery"`
	Uptime           int                `bson:"uptime"`
	CoreTemp         int                `bson:"core_temp"`
	ReceivedOn       string             `bson:"received_on"`
}

// Telemetry raw telemetry data
type Telemetry struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
}

// TelemetryEvent for storing drone telemetry data
type TelemetryEvent struct {
	Telemetry  Telemetry `json:"telemetry"`
	ReceivedOn int64     `json:"received_on"`
}

// AlertRecord for storing drone failure alert data
type AlertRecord struct {
	RecordID    primitive.ObjectID `bson:"_id,omitempty"`
	DroneID     string             `bson:"drone_id"`
	FaultCode   int                `bson:"fault_code"`
	Description string             `bson:"description"`
	ReceivedOn  string             `bson:"received_on"`
}

// Alert raw alert data
type Alert struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
}

// AlertEvent for receiving drone failure alert data
type AlertEvent struct {
	ReceivedOn int64 `json:"received_on"`
	Alert      Alert `json:"alert"`
}

// PositionRecord for submitting drone position updates data
type PositionRecord struct {
	RecordID        primitive.ObjectID `bson:"_id,omitempty"`
	DroneID         string             `bson:"drone_id"`
	Latitude        float32            `bson:"latitude"`
	Longitude       float32            `bson:"longitude"`
	Altitude        float32            `bson:"altitude"`
	CurrentSpeed    float32            `bson:"current_speed"`
	HeadingCardinal int                `bson:"heading_cardinal"`
	ReceivedOn      string             `bson:"received_on"`
}

// Position raw position data
type Position struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
}

// PositionEvent for submitting drone position updates data
type PositionEvent struct {
	Position   Position `json:"position"`
	ReceivedOn int64    `json:"received_on"`
}

func convertToTelemetryRecord(event TelemetryEvent) TelemetryRecord {
	t := time.Unix(event.ReceivedOn, 0)
	return TelemetryRecord{
		DroneID:          event.Telemetry.DroneID,
		RemainingBattery: event.Telemetry.RemainingBattery,
		Uptime:           event.Telemetry.Uptime,
		CoreTemp:         event.Telemetry.CoreTemp,
		ReceivedOn:       t.Format(time.UnixDate),
	}
}

func convertToAlertRecord(event AlertEvent) AlertRecord {
	t := time.Unix(event.ReceivedOn, 0)
	return AlertRecord{
		DroneID:     event.Alert.DroneID,
		Description: event.Alert.Description,
		FaultCode:   event.Alert.FaultCode,
		ReceivedOn:  t.Format(time.UnixDate),
	}
}

func convertToPositionRecord(event PositionEvent) PositionRecord {
	t := time.Unix(event.ReceivedOn, 0)
	return PositionRecord{
		DroneID:         event.Position.DroneID,
		Latitude:        event.Position.Latitude,
		Longitude:       event.Position.Longitude,
		Altitude:        event.Position.Altitude,
		CurrentSpeed:    event.Position.CurrentSpeed,
		HeadingCardinal: event.Position.HeadingCardinal,
		ReceivedOn:      t.Format(time.UnixDate),
	}
}
