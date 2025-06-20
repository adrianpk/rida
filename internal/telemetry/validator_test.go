package telemetry_test

import (
	"errors"
	"testing"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

func TestDefaultValidator(t *testing.T) {
	validID := uuid.New()
	validArea := telemetry.Area{MinLat: 1, MaxLat: 2, MinLng: 3, MaxLng: 4}
	invalidArea := telemetry.Area{MinLat: 2, MaxLat: 1, MinLng: 4, MaxLng: 3}

	tests := []struct {
		name    string
		op      telemetry.ValidationOp
		data    any
		wantErr error
	}{
		{
			name:    "valid get scooter",
			op:      telemetry.OpGetScooter,
			data:    validID,
			wantErr: nil,
		},
		{
			name:    "invalid get scooter (nil id)",
			op:      telemetry.OpGetScooter,
			data:    uuid.Nil,
			wantErr: telemetry.ErrInvalidID,
		},
		{
			name:    "valid update scooter",
			op:      telemetry.OpUpdateScooter,
			data:    telemetry.Scooter{ID: validID},
			wantErr: nil,
		},
		{
			name:    "invalid update scooter (nil id)",
			op:      telemetry.OpUpdateScooter,
			data:    telemetry.Scooter{ID: uuid.Nil},
			wantErr: telemetry.ErrInvalidID,
		},
		{
			name:    "valid find scooters",
			op:      telemetry.OpFindScooters,
			data:    telemetry.Query{Area: validArea, Status: "free"},
			wantErr: nil,
		},
		{
			name:    "invalid find scooters (area bounds)",
			op:      telemetry.OpFindScooters,
			data:    telemetry.Query{Area: invalidArea, Status: "free"},
			wantErr: errors.New("invalid area bounds"),
		},
		{
			name:    "invalid find scooters (wrong type)",
			op:      telemetry.OpFindScooters,
			data:    123,
			wantErr: errors.New("invalid query params"),
		},
		{
			name:    "invalid report event (bad type)",
			op:      telemetry.OpReportEvent,
			data:    telemetry.Event{ID: validID, ScooterID: validID, Type: "bad_type"},
			wantErr: errors.New("invalid event type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := telemetry.DefaultValidator(tt.op, tt.data)
			if (err == nil) != (tt.wantErr == nil) {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestIsValidEventType(t *testing.T) {
	tests := []struct {
		name  string
		input telemetry.EventType
		want  bool
	}{
		{"trip_start valid", telemetry.EventTripStart, true},
		{"trip_end valid", telemetry.EventTripEnd, true},
		{"location valid", telemetry.EventLocation, true},
		{"invalid type", "bad_type", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := telemetry.IsValidEventType(tt.input)
			if got != tt.want {
				t.Errorf("IsValidEventType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
