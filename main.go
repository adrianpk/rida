package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianpk/rida/internal/cfg"
	"github.com/adrianpk/rida/internal/client"
	"github.com/adrianpk/rida/internal/repo/pg"
	"github.com/adrianpk/rida/internal/telemetry"
)

const (
	AppName    = "Rida"
	AppVersion = "1.0.0"
)

func main() {
	config := cfg.Load()

	db := pg.NewDB(config)
	repo := pg.NewTelemetryRepo(db)

	err := startDeps(db, repo)
	if err != nil {
		log.Fatal(err)
	}

	service := telemetry.NewService(repo)
	handler := telemetry.NewHandler(service)
	router := telemetry.NewRouter(handler, config.APIKey)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startClients(ctx, config)
	startServer(router, config)
}

func startDeps(db *pg.DB, repo *pg.TelemetryRepo) error {
	ctx := context.Background()
	if err := db.Setup(ctx); err != nil {
		return err
	}

	if err := repo.Setup(ctx); err != nil {
		return err
	}

	return nil
}

func startServer(router http.Handler, config *cfg.Config) {
	log.Printf("%s running on %s", AppName, config.HTTPPort)
	log.Fatal(http.ListenAndServe(config.HTTPPort, router))
}

func startClients(ctx context.Context, config *cfg.Config) {
	manager := client.NewClientManager(config)
	manager.Start(ctx)
}
