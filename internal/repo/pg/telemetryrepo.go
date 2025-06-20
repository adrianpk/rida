package pg

import (
	"context"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

// TelemetryRepo is a PostgreSQL implementation of the telemetry.Repo interface.
type TelemetryRepo struct {
	db *DB
}

// NewTelemetryRepo creates a new PostgreSQL-backed TelemetryRepo.
func NewTelemetryRepo(db *DB) *TelemetryRepo {
	return &TelemetryRepo{db: db}
}

func (r *TelemetryRepo) GetScooter(ctx context.Context, id uuid.UUID) (telemetry.Scooter, error) {
	var scooter telemetry.Scooter
	q := query[getScooterQueryKey]
	err := r.db.GetContext(ctx, &scooter, q, id)

	return scooter, err
}

func (r *TelemetryRepo) UpdateScooter(ctx context.Context, s telemetry.Scooter) error {
	q := query[updateScooterQueryKey]
	_, err := r.db.NamedExecContext(ctx, q, s)

	return err
}

func (r *TelemetryRepo) FindScootersInArea(ctx context.Context, area telemetry.Area, status telemetry.Status) ([]telemetry.Scooter, error) {
	q := query[findScootersInAreaQueryKey]
	var scooters []telemetry.Scooter
	rows, err := r.db.NamedQueryContext(ctx, q, map[string]interface{}{
		"min_lat": area.MinLat,
		"max_lat": area.MaxLat,
		"min_lng": area.MinLng,
		"max_lng": area.MaxLng,
		"status":  status,
	})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s telemetry.Scooter
		if err := rows.StructScan(&s); err != nil {
			return nil, err
		}
		scooters = append(scooters, s)
	}

	return scooters, nil
}

func (r *TelemetryRepo) StoreEvent(ctx context.Context, e telemetry.Event) error {
	q := query[storeEventQueryKey]
	_, err := r.db.NamedExecContext(ctx, q, e)

	return err
}

// Setup runs migration and seeding for the TelemetryRepo.
func (r *TelemetryRepo) Setup(ctx context.Context) error {
	err := r.Migrate(ctx)
	if err != nil {
		return err
	}

	err = r.Seed(ctx)
	if err != nil {
		return err
	}

	return nil
}
