package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianpk/rida/internal/cfg"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB wraps sqlx.DB to allow custom methods such as Setup, and to enable future
// extensions like logging, observability, error tracing, and other
// cross-cutting concerns.
type DB struct {
	cfg *cfg.Config
	*sqlx.DB
}

func NewDB(cfg *cfg.Config) *DB {
	return &DB{DB: sqlx.NewDb(nil, "postgres"), cfg: cfg}
}

func (db *DB) Setup(ctx context.Context) error {
	var pgdb *sqlx.DB
	var err error

	for i := 0; i < 10; i++ {
		pgdb, err = sqlx.ConnectContext(ctx, "postgres", db.cfg.Pg.DSN())
		if err == nil {
			db.DB = pgdb
			return nil
		}

		waitTime := time.Duration(i+1) * time.Second
		if waitTime > 10*time.Second {
			waitTime = 10 * time.Second
		}

		fmt.Printf("postgres connection attempt %d failed, retrying in %v: %v\n", i+1, waitTime, err)
		time.Sleep(waitTime)
	}

	return fmt.Errorf("postgres connection failed after 10 attempts: %w", err)
}
