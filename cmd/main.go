package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"okta-ingestor/internal/config"
	"okta-ingestor/internal/okta"
	"okta-ingestor/internal/scheduler"
	"okta-ingestor/internal/server"
	"okta-ingestor/internal/storage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, _ := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())

	server.StartHealthServer()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI)) // db connection
	db := client.Database(cfg.Mongo.Database)                                 // connection to make db comms

	repo := storage.NewMongoRepository(db, cfg.Mongo.LogsCol, cfg.Mongo.StateCol)
	hook := storage.NewWebhookSender(cfg.Webhook.URL, cfg.Webhook.Method, cfg.Webhook.Timeout, cfg.Webhook.Retries)
	oktaClient := okta.New(cfg.Okta.Domain, cfg.Okta.Token, cfg.Okta.Limit)

	// initialSince := time.Now().Add(-1 * time.Hour)
	initialSince, _ := time.Parse(time.RFC3339, cfg.Okta.Since)

	s := scheduler.New(oktaClient, repo, hook, cfg.PollInterval)
	go s.Start(ctx, initialSince)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop   //
	cancel() // cancel the scheduler context to stop it gracefully
	slog.Info("shutdown complete")
}
