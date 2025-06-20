package telemetry

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetScooter(ctx context.Context, id uuid.UUID) (Scooter, error)
	UpdateScooter(ctx context.Context, s Scooter) error
	FindScooters(ctx context.Context, qry Query) ([]Scooter, error)
	ReportEvent(ctx context.Context, e Event) error
}

type service struct {
	repo     Repo
	validate Validator
}

func NewService(r Repo) Service {
	return &service{
		repo:     r,
		validate: DefaultValidator,
	}
}

func (s *service) GetScooter(ctx context.Context, id uuid.UUID) (Scooter, error) {
	err := s.validate(OpGetScooter, id)
	if err != nil {
		return Scooter{}, err
	}

	return s.repo.GetScooter(ctx, id)
}

func (s *service) UpdateScooter(ctx context.Context, scooter Scooter) error {
	err := s.validate(OpUpdateScooter, scooter)
	if err != nil {
		return err
	}

	return s.repo.UpdateScooter(ctx, scooter)
}

func (s *service) FindScooters(ctx context.Context, qry Query) ([]Scooter, error) {
	err := s.validate(OpFindScooters, qry)
	if err != nil {
		return nil, err
	}

	return s.repo.FindScootersInArea(ctx, qry.Area, qry.Status)
}

// ReportEvent processes an incoming event and updates the scooter state accordingly.
//
// NOTE: In a production system, an event streaming approach (e.g., using NATS)
// could be used for decoupling, scalability, and reliability. For this home assignment,
// we use a simpler approach: events are processed synchronously and
// directly update the scooter state if no errors occur. See docs/adr/0001-event-processing-vs-streaming.md.
func (s *service) ReportEvent(ctx context.Context, e Event) error {
	if err := s.validate(OpReportEvent, e); err != nil {
		return err
	}

	e.GenCreateVals()

	err := s.repo.StoreEvent(ctx, e)
	if err != nil {
		return err
	}

	scooter, err := s.repo.GetScooter(ctx, e.ScooterID)
	if err != nil {
		return err
	}

	switch e.Type {
	case EventTripStart:
		scooter.StartRide()
	case EventTripEnd:
		scooter.StopRide()
	case EventLocation:
		scooter.UpdateLocation(e.Lat, e.Lng)
	}

	return s.repo.UpdateScooter(ctx, scooter)
}

// SetValidator lets you replace the default validator with a custom one.
// This is handy in tests to inject mock or specialized validation logic.
func (s *service) SetValidator(v Validator) {
	s.validate = v
}
