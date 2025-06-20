package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

// TelemetryRepo is an in-memory implementation of the Repo interface.
//
// This implementation is intended for development and testing purposes. In a production
// environment, a robust persistence layer (such as a SQL or NoSQL database) should be used.
//
// Transactions: The repository interface does not expose specific methods for transaction
// management. Each concrete implementation should handle transactions as appropriate,
// ideally obtaining them from the context. This keeps the interface simple and decoupled
// from infrastructure details, making it easier to integrate with different database engines
// and persistence patterns.
type TelemetryRepo struct {
	mu       sync.RWMutex
	scooters map[uuid.UUID]telemetry.Scooter
	events   []telemetry.Event
}

func NewTelemetryRepo(initial ...map[uuid.UUID]telemetry.Scooter) *TelemetryRepo {
	var scooters map[uuid.UUID]telemetry.Scooter
	if len(initial) > 0 {
		scooters = initial[0]
	}

	if scooters == nil {
		scooters = make(map[uuid.UUID]telemetry.Scooter)
	}

	repo := &TelemetryRepo{
		scooters: scooters,
	}

	return repo
}

func (r *TelemetryRepo) GetScooter(ctx context.Context, id uuid.UUID) (telemetry.Scooter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.scooters[id]
	if !ok {
		return telemetry.Scooter{}, errors.New("not found")
	}

	return s, nil
}

func (r *TelemetryRepo) UpdateScooter(ctx context.Context, s telemetry.Scooter) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scooters[s.ID] = s
	return nil
}

func (r *TelemetryRepo) FindScootersInArea(ctx context.Context, area telemetry.Area, status telemetry.Status) ([]telemetry.Scooter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []telemetry.Scooter
	for _, s := range r.scooters {
		if s.Status == status &&
			s.Lat >= area.MinLat && s.Lat <= area.MaxLat &&
			s.Lng >= area.MinLng && s.Lng <= area.MaxLng {
			result = append(result, s)
		}
	}

	return result, nil
}

func (r *TelemetryRepo) StoreEvent(ctx context.Context, e telemetry.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
	return nil
}

// Scooters returns a copy of the scooters map for black-box testing.
func (r *TelemetryRepo) Scooters() map[uuid.UUID]telemetry.Scooter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copyMap := make(map[uuid.UUID]telemetry.Scooter, len(r.scooters))
	for k, v := range r.scooters {
		copyMap[k] = v
	}
	return copyMap
}

// Events returns a copy of the events slice for black-box testing.
func (r *TelemetryRepo) Events() []telemetry.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]telemetry.Event(nil), r.events...)
}
