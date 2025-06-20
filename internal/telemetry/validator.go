package telemetry

import (
	"errors"

	"github.com/google/uuid"
)

type ValidationOp string

const (
	OpGetScooter    ValidationOp = "get"
	OpUpdateScooter ValidationOp = "update"
	OpFindScooters  ValidationOp = "find"
	OpReportEvent   ValidationOp = "report_event"
)

type Validator func(op ValidationOp, data interface{}) error

var ErrInvalidID = errors.New("invalid scooter id")

func DefaultValidator(op ValidationOp, data interface{}) error {
	switch op {
	case OpGetScooter:
		id, ok := data.(uuid.UUID)
		if !ok || id == uuid.Nil {
			return ErrInvalidID
		}

	case OpUpdateScooter:
		scooter, ok := data.(Scooter)
		if !ok || scooter.ID == uuid.Nil {
			return ErrInvalidID
		}

	case OpFindScooters:
		params, ok := data.(Query)
		if !ok {
			return errors.New("invalid query params")
		}

		if params.Area.MinLat > params.Area.MaxLat ||
			params.Area.MinLng > params.Area.MaxLng {
			return errors.New("invalid area bounds")
		}

	case OpReportEvent:
		return validateReportEvent(data)
	}

	return nil
}

// OpReportEvent validates an event for reporting.
func validateReportEvent(v interface{}) error {
	e, ok := v.(Event)
	if !ok {
		return errors.New("invalid event type")
	}

	if e.ScooterID == uuid.Nil {
		return errors.New("invalid scooter id")
	}

	if !IsValidEventType(e.Type) {
		return errors.New("invalid event type")
	}

	return nil
}

func IsValidEventType(t EventType) bool {
	switch t {
	case EventTripStart, EventTripEnd, EventLocation:
		return true

	default:
		return false
	}
}
