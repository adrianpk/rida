package telemetry_test

import (
	"testing"
	"time"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

func TestScooterStartRide(t *testing.T) {
	s := &telemetry.Scooter{
		ID:        uuid.New(),
		Status:    telemetry.StatusFree,
		Lat:       0,
		Lng:       0,
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	previous := s.UpdatedAt

	s.StartRide()

	if s.Status != telemetry.StatusOccupied {
		t.Errorf("expected status %q, got %q", telemetry.StatusOccupied, s.Status)
	}

	if !s.UpdatedAt.After(previous) {
		t.Errorf("expected updated at to be after previous, got %v", s.UpdatedAt)
	}
}

func TestScooterStopRide(t *testing.T) {
	s := &telemetry.Scooter{
		ID:        uuid.New(),
		Status:    telemetry.StatusOccupied,
		Lat:       0,
		Lng:       0,
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	previous := s.UpdatedAt

	s.StopRide()

	if s.Status != telemetry.StatusFree {
		t.Errorf("expected status %q, got %q", telemetry.StatusFree, s.Status)
	}

	if !s.UpdatedAt.After(previous) {
		t.Errorf("expected updated at to be after previous, got %v", s.UpdatedAt)
	}
}

func TestScooterGenID(t *testing.T) {
	tests := []struct {
		name     string
		input    telemetry.Scooter
		wantZero bool
	}{
		{
			name:     "already has ID",
			input:    telemetry.Scooter{ID: uuid.New()},
			wantZero: false,
		},
		{
			name:     "no ID set",
			input:    telemetry.Scooter{},
			wantZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			s.GenID()
			if s.ID == uuid.Nil {
				t.Errorf("expected non-nil ID after GenID, got nil")
			}
		})
	}
}
