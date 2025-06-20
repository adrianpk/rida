package telemetry

import (
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	GetScooter(ctx context.Context, id uuid.UUID) (Scooter, error)
	UpdateScooter(ctx context.Context, s Scooter) error
	FindScootersInArea(ctx context.Context, area Area, status Status) ([]Scooter, error)
	StoreEvent(ctx context.Context, e Event) error
}
