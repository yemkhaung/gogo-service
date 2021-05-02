package dronescmds

type telemetryCommand struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
}

type alertCommand struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
}

type positionCommand struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
}

func (t telemetryCommand) isValid() bool {
	if t.DroneID == "" {
		return false
	}
	if t.Uptime == 0 {
		return false
	}
	return true
}

func (a alertCommand) isValid() bool {
	if a.DroneID == "" {
		return false
	}
	if a.FaultCode == 0 {
		return false
	}
	return true
}

func (p positionCommand) isValid() bool {
	if p.DroneID == "" {
		return false
	}
	if p.Longitude < 0 || p.Latitude < 0 || p.HeadingCardinal < 0 {
		return false
	}
	return true
}

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
