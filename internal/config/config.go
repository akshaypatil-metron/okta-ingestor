package config

import "time"

type Config struct {
	OktaDomain   string
	OktaToken    string
	MongoURI     string
	WebhookURL   string
	PollInterval time.Duration // Only ONE timer needed now
	BatchSize    int
	StartDate    string
}

func Load() *Config {
	return &Config{
		OktaDomain: "https://trial-5020092.okta.com",
		OktaToken:  "00M2hSXzGtbugWg-5YYPksVH4o7MjtxEQ-gcNLJA1t",
		MongoURI:   "mongodb://localhost:27017",
		WebhookURL: "https://my-okta-logs.free.beeceptor.com",

		PollInterval: 30 * time.Second, // Pause for 30s between EVERY fetch
		BatchSize:    100,              // Fetch 100 at a time
		StartDate:    "2026-02-01T00:00:00.000Z",
	}
}
