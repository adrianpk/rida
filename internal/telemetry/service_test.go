package telemetry_test

import (
	"context"
	"testing"

	"github.com/adrianpk/rida/internal/repo/mem"
	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

func TestService_ReportEvent(t *testing.T) {
	type args struct {
		event telemetry.Event
	}

	scooterID := uuid.New()
	initialScooter := telemetry.Scooter{
		ID:     scooterID,
		Status: telemetry.StatusFree,
		Lat:    45.0,
		Lng:    -75.0,
	}

	tests := []struct {
		name       string
		initial    map[uuid.UUID]telemetry.Scooter
		args       args
		wantStatus telemetry.Status
		wantLat    float64
		wantLng    float64
		wantErr    bool
	}{
		{
			name:    "trip start updates status to occupied",
			initial: initialData(initialScooter),
			args: args{
				event: telemetry.Event{
					ScooterID: scooterID,
					Type:      telemetry.EventTripStart,
				},
			},
			wantStatus: telemetry.StatusOccupied,
			wantLat:    initialScooter.Lat,
			wantLng:    initialScooter.Lng,
			wantErr:    false,
		},
		{
			name:    "trip end updates status to free",
			initial: initialData(initialScooter),
			args: args{
				event: telemetry.Event{
					ScooterID: scooterID,
					Type:      telemetry.EventTripEnd,
				},
			},
			wantStatus: telemetry.StatusFree,
			wantLat:    initialScooter.Lat,
			wantLng:    initialScooter.Lng,
			wantErr:    false,
		},
		{
			name:    "location update changes lat/lng",
			initial: initialData(initialScooter),
			args: args{
				event: telemetry.Event{
					ScooterID: scooterID,
					Type:      telemetry.EventLocation,
					Lat:       46.0,
					Lng:       -76.0,
				},
			},
			wantStatus: telemetry.StatusFree,
			wantLat:    46.0,
			wantLng:    -76.0,
			wantErr:    false,
		},
		{
			name:    "invalid event type returns error",
			initial: initialData(initialScooter),
			args: args{
				event: telemetry.Event{
					ScooterID: scooterID,
					Type:      "bad_type",
				},
			},
			wantStatus: initialScooter.Status,
			wantLat:    initialScooter.Lat,
			wantLng:    initialScooter.Lng,
			wantErr:    true,
		},
		{
			name:    "nonexistent scooter returns error",
			initial: map[uuid.UUID]telemetry.Scooter{},
			args: args{
				event: telemetry.Event{
					ScooterID: uuid.New(),
					Type:      telemetry.EventTripStart,
				},
			},
			wantStatus: "",
			wantLat:    0,
			wantLng:    0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mem.NewTelemetryRepo(tt.initial)
			svc := telemetry.NewService(repo)

			var scooterID uuid.UUID
			scooters := repo.Scooters()
			for id := range scooters {
				scooterID = id
				break
			}

			if scooterID != uuid.Nil && tt.args.event.ScooterID == uuid.Nil {
				tt.args.event.ScooterID = scooterID
			}

			err := svc.ReportEvent(context.Background(), tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReportEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			scooters = repo.Scooters()

			s, ok := scooters[tt.args.event.ScooterID]
			if !ok {
				t.Fatalf("Scooter not found after event processing")
			}

			if s.Status != tt.wantStatus {
				t.Errorf("Scooter status = %v, want %v", s.Status, tt.wantStatus)
			}

			if s.Lat != tt.wantLat || s.Lng != tt.wantLng {
				t.Errorf("Scooter location = (%v, %v), want (%v, %v)", s.Lat, s.Lng, tt.wantLat, tt.wantLng)
			}
		})
	}
}

func initialData(scooter telemetry.Scooter) map[uuid.UUID]telemetry.Scooter {
	return map[uuid.UUID]telemetry.Scooter{scooter.ID: scooter}
}
