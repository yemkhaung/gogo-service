package dronesevents

type eventRepository interface {
	UpdateLastTelemetryEvent(TelemetryUpdatedEvent) error
	UpdateLastAlertEvent(AlertSignalledEvent) error
	UpdateLastPositionEvent(PositionChangedEvent) error
}

type mongoEventRepository struct {
	URL string
}

func (repo *mongoEventRepository) UpdateLastTelemetryEvent(event TelemetryUpdatedEvent) error {

	return nil
}

func (repo *mongoEventRepository) UpdateLastAlertEvent(event AlertSignalledEvent) error {

	return nil
}

func (repo *mongoEventRepository) UpdateLastPositionEvent(event PositionChangedEvent) error {

	return nil
}
