package mem_test

import (
	"context"
	"testing"

	"github.com/adrianpk/rida/internal/repo/mem"
	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

func TestGetScooterByID(t *testing.T) {
	repo := mem.NewTelemetryRepo()
	testScooter := telemetry.Scooter{
		ID:     uuid.New(),
		Status: telemetry.StatusFree,
		Lat:    51.1079,
		Lng:    17.0385,
	}

	err := repo.UpdateScooter(context.Background(), testScooter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "found",
			id:      testScooter.ID,
			wantErr: false,
		},
		{
			name:    "not found",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetScooter(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdateScooter(t *testing.T) {
	scooterID := uuid.New()
	repo := mem.NewTelemetryRepo()

	tests := []struct {
		name      string
		initial   map[uuid.UUID]telemetry.Scooter
		update    telemetry.Scooter
		wantErr   bool
		checkFunc func(repo *mem.TelemetryRepo, s telemetry.Scooter) bool
	}{
		{
			name: "update existing scooter",
			initial: map[uuid.UUID]telemetry.Scooter{
				scooterID: {
					ID:     scooterID,
					Status: telemetry.StatusFree,
				},
			},
			update: telemetry.Scooter{
				ID:     scooterID,
				Status: telemetry.StatusOccupied,
			},
			wantErr: false,
			checkFunc: func(repo *mem.TelemetryRepo, s telemetry.Scooter) bool {
				scooters := repo.Scooters()
				stored, ok := scooters[s.ID]
				return ok && stored.Status == s.Status
			},
		},
		{
			name:    "add new scooter",
			initial: map[uuid.UUID]telemetry.Scooter{},
			update: telemetry.Scooter{
				ID:     uuid.New(),
				Status: telemetry.StatusFree,
			},
			wantErr: false,
			checkFunc: func(repo *mem.TelemetryRepo, s telemetry.Scooter) bool {
				scooters := repo.Scooters()
				stored, ok := scooters[s.ID]
				return ok && stored.Status == s.Status
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo = mem.NewTelemetryRepo(tt.initial)
			err := repo.UpdateScooter(context.Background(), tt.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !tt.checkFunc(repo, tt.update) {
				t.Errorf("scooter not updated as expected")
			}
		})
	}
}

func TestFindScootersInArea(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	id4 := uuid.New()

	tests := []struct {
		name    string
		initial map[uuid.UUID]telemetry.Scooter
		area    telemetry.Area
		status  telemetry.Status
		wantIDs []uuid.UUID
	}{
		{
			name: "one inside area and status",
			initial: map[uuid.UUID]telemetry.Scooter{
				id1: {ID: id1, Lat: 51.5, Lng: 17.5, Status: telemetry.StatusFree},     // inside
				id2: {ID: id2, Lat: 52.1, Lng: 17.5, Status: telemetry.StatusFree},     // outside (lat)
				id3: {ID: id3, Lat: 51.5, Lng: 18.1, Status: telemetry.StatusFree},     // outside (lng)
				id4: {ID: id4, Lat: 51.5, Lng: 17.5, Status: telemetry.StatusOccupied}, // inside, wrong status
			},
			area: telemetry.Area{
				MinLat: 51.0, MaxLat: 52.0,
				MinLng: 17.0, MaxLng: 18.0,
			},
			status:  telemetry.StatusFree,
			wantIDs: []uuid.UUID{id1},
		},
		{
			name: "none in area",
			initial: map[uuid.UUID]telemetry.Scooter{
				id2: {ID: id2, Lat: 52.1, Lng: 17.5, Status: telemetry.StatusFree},
				id3: {ID: id3, Lat: 51.5, Lng: 18.1, Status: telemetry.StatusFree},
			},
			area: telemetry.Area{
				MinLat: 51.0, MaxLat: 52.0,
				MinLng: 17.0, MaxLng: 18.0,
			},
			status:  telemetry.StatusOccupied,
			wantIDs: nil,
		},
		{
			name: "multiple inside area and status",
			initial: map[uuid.UUID]telemetry.Scooter{
				id1: {ID: id1, Lat: 51.1, Lng: 17.1, Status: telemetry.StatusFree},
				id2: {ID: id2, Lat: 51.2, Lng: 17.2, Status: telemetry.StatusFree},
				id3: {ID: id3, Lat: 51.3, Lng: 17.3, Status: telemetry.StatusOccupied},
			},
			area: telemetry.Area{
				MinLat: 51.0, MaxLat: 51.25,
				MinLng: 17.0, MaxLng: 17.25,
			},
			status:  telemetry.StatusFree,
			wantIDs: []uuid.UUID{id1, id2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mem.NewTelemetryRepo(tt.initial)
			got, err := repo.FindScootersInArea(context.Background(), tt.area, tt.status)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gotIDs []uuid.UUID
			for _, s := range got {
				gotIDs = append(gotIDs, s.ID)
			}

			if !equalUUIDSlices(gotIDs, tt.wantIDs) {
				t.Errorf("got IDs %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func equalUUIDSlices(a, b []uuid.UUID) bool {
	if len(a) != len(b) {
		return false
	}

	m := make(map[uuid.UUID]int)
	for _, id := range a {
		m[id]++
	}

	for _, id := range b {
		if m[id] == 0 {
			return false
		}
		m[id]--
	}

	return true
}

func TestStoreEvent(t *testing.T) {
	repo := mem.NewTelemetryRepo()
	events := []telemetry.Event{
		{
			ID:        uuid.New(),
			ScooterID: uuid.New(),
			Type:      telemetry.EventTripStart,
			Lat:       51.1,
			Lng:       17.0,
		},
		{
			ID:        uuid.New(),
			ScooterID: uuid.New(),
			Type:      telemetry.EventTripEnd,
			Lat:       51.2,
			Lng:       17.1,
		},
		{
			ID:        uuid.New(),
			ScooterID: uuid.New(),
			Type:      telemetry.EventLocation,
			Lat:       51.3,
			Lng:       17.2,
		},
	}

	for i, event := range events {
		err := repo.StoreEvent(context.Background(), event)
		if err != nil {
			t.Fatalf("event %d: expected no error, got %v", i, err)
		}
	}

	stored := repo.Events()
	if len(stored) != len(events) {
		t.Fatalf("expected %d events, got %d", len(events), len(stored))
	}
	for i, event := range events {
		storedEvent := stored[i]
		if storedEvent != event {
			t.Errorf("event %d does not match: got %+v, want %+v", i, storedEvent, event)
		}
	}
}
