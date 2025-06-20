package pg

import (
	"context"
	"math/rand"

	"github.com/adrianpk/rida/internal/telemetry"
)

type CitySeed struct {
	Name  string
	Count int
	Area  telemetry.Area
}

// Seed inserts demo scooters for Ottawa and Montreal into the PostgreSQL database.
func (r *TelemetryRepo) Seed(ctx context.Context) error {
	cities := []CitySeed{
		{
			Name:  "Ottawa",
			Count: 3216,
			Area: telemetry.Area{
				MinLat: 45.17927019403111,
				MaxLat: 45.4502599310963,
				MinLng: -75.95781905735376,
				MaxLng: -75.37765015636133,
			},
		},
		{
			Name:  "Montreal",
			Count: 5376,
			Area: telemetry.Area{
				MinLat: 45.452507945877,
				MaxLat: 45.62109228798646,
				MinLng: -73.63465335011105,
				MaxLng: -73.55019903119938,
			},
		},
	}

	for _, city := range cities {
		if err := r.SeedCity(ctx, city.Count, city.Area); err != nil {
			return err
		}
	}
	return nil
}

// SeedCity inserts a specific number of scooters randomly in the given area.
func (r *TelemetryRepo) SeedCity(ctx context.Context, count int, area telemetry.Area) error {
	for i := 0; i < count; i++ {
		status := telemetry.StatusFree
		if rand.Float64() < 0.6 {
			status = telemetry.StatusOccupied
		}
		s := telemetry.Scooter{
			Status: status,
			Lat:    area.MinLat + rand.Float64()*(area.MaxLat-area.MinLat),
			Lng:    area.MinLng + rand.Float64()*(area.MaxLng-area.MinLng),
		}
		s.GenID()
		// Insert scooter into DB
		_, err := r.db.ExecContext(ctx,
			`INSERT INTO scooters (id, status, lat, lng, updated_at) VALUES ($1, $2, $3, $4, $5)`,
			s.ID, s.Status, s.Lat, s.Lng, s.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
