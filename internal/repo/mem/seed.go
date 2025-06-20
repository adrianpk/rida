package mem

import (
	"math/rand"

	"github.com/adrianpk/rida/internal/telemetry"
)

type CitySeed struct {
	Name  string
	Count int
	Area  telemetry.Area
}

// Seed seeds the default set of cities with scooters for demo
// or testing purposes.
func (r *TelemetryRepo) Seed() {
	r.SeedCities()
}

func (r *TelemetryRepo) SeedCities() {
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
		r.SeedCity(city.Count, city.Area)
	}
}

// SeedCity seeds a specific number of scooters randomly in the given area.
func (r *TelemetryRepo) SeedCity(count int, area telemetry.Area) {
	r.mu.Lock()
	defer r.mu.Unlock()
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

		r.scooters[s.ID] = s
	}
}
