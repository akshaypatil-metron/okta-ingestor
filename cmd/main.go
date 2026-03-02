package main

import (
	"log"

	"okta-ingestor/internal/config"
	"okta-ingestor/internal/okta"
	"okta-ingestor/internal/processor"
	"okta-ingestor/internal/scheduler"
	"okta-ingestor/internal/server"
	"okta-ingestor/internal/storage"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Infrastructure / Clients
	mongoStore, err := storage.NewMongoStore(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	oktaClient := okta.NewClient(cfg.OktaToken)
	webhookClient := storage.NewWebhookClient(cfg.WebhookURL)

	// 3. Initialize Processor
	proc := processor.NewProcessor(oktaClient, mongoStore, webhookClient, cfg.OktaDomain, cfg.BatchSize, cfg.StartDate)

	// 4. Start Health Check Server (Non-blocking)
	go server.StartHealthCheck()

	// 5. Start Scheduler (Blocking)
	sched := scheduler.NewScheduler(proc, cfg.PollInterval)
	sched.Start()
}
