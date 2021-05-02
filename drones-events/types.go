package dronesevents

// TelemetryCommand for submitting drone telemetry data
type TelemetryCommand struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
}

// AlertCommand for submitting drone failure alert data
type AlertCommand struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
}

// PositionCommand for submitting drone position updates data
type PositionCommand struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
}

// TelemetryUpdatedEvent is an event containing telemetry updates received from a drone
type TelemetryUpdatedEvent struct {
	Telemetry  TelemetryCommand `json:"telemetry"`
	ReceivedOn int64            `json:"received_on"`
}

// AlertSignalledEvent is an event indicating an alert condition was reported by a drone
type AlertSignalledEvent struct {
	Alert      AlertCommand `json:"alert"`
	ReceivedOn int64        `json:"received_on"`
}

// PositionChangedEvent is an event indicating that the position and speed of a drone was reported.
type PositionChangedEvent struct {
	Position   PositionCommand `json:"position"`
	ReceivedOn int64           `json:"received_on"`
}
