package telemetry

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusFree     Status = "free"
	StatusOccupied Status = "occupied"
)

type Scooter struct {
	ID        uuid.UUID `json:"id"`
	Status    Status    `json:"status"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Scooter) GenID() {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
}

func (s *Scooter) StartRide() {
	s.Status = StatusOccupied
	s.AuditUpdate()
}

func (s *Scooter) StopRide() {
	s.Status = StatusFree
	s.AuditUpdate()
}

func (s *Scooter) UpdateLocation(lat, lng float64) {
	s.Lat = lat
	s.Lng = lng
	s.AuditUpdate()
}

func (s *Scooter) AuditUpdate() {
	s.UpdatedAt = time.Now()
}

func (s *Scooter) GenCreateVals() {
	s.GenID()
	s.UpdatedAt = time.Now()
}

type EventType string

const (
	EventTripStart EventType = "trip_start"
	EventTripEnd   EventType = "trip_end"
	EventLocation  EventType = "location"
)

type Event struct {
	ID        uuid.UUID `json:"id"`
	ScooterID uuid.UUID `json:"scooterId"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
}

func (e *Event) GenID() {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
}

func (e *Event) GenCreateVals() {
	e.GenID()
	e.Timestamp = time.Now()
}

type Area struct {
	MinLat float64 `json:"minLat"`
	MinLng float64 `json:"minLng"`
	MaxLat float64 `json:"maxLat"`
	MaxLng float64 `json:"maxLng"`
}

type Query struct {
	Area   Area
	Status Status
}
