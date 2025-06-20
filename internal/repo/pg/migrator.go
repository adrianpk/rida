package pg

import (
	"context"
	"fmt"
)

// Migrate creates the tables needed for Scooter and Event in a simple way.
// This is a basic implementation just to satisfy the use case for this project.
func (r *TelemetryRepo) Migrate(ctx context.Context) error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS postgis;`,
		`CREATE TABLE IF NOT EXISTS scooters (
			id UUID PRIMARY KEY,
			status TEXT NOT NULL,
			lat DOUBLE PRECISION NOT NULL,
			lng DOUBLE PRECISION NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			scooter_id UUID NOT NULL,
			type TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			lat DOUBLE PRECISION NOT NULL,
			lng DOUBLE PRECISION NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
