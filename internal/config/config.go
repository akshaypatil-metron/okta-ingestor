package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Okta struct {
		Domain string
		Token  string
		Since  string
		Limit  int
	}

	PollInterval time.Duration

	Mongo struct {
		URI      string
		Database string
		LogsCol  string
		StateCol string
	}

	Webhook struct {
		URL     string
		Method  string
		Timeout time.Duration
		Retries int
	}
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("OKTA_LIMIT", 10)
	viper.SetDefault("POLL_INTERVAL", "30s")
	viper.SetDefault("MONGO_LOGS_COL", "logs")
	viper.SetDefault("MONGO_STATE_COL", "state")
	viper.SetDefault("WEBHOOK_METHOD", "POST")
	viper.SetDefault("WEBHOOK_TIMEOUT", "10s")
	viper.SetDefault("WEBHOOK_RETRIES", 3)
	viper.SetDefault("OKTA_DOMAIN", "trial-4733557.okta.com")
	viper.SetDefault("OKTA_TOKEN", "00rcAZ_SB-2rKAcS9druJJv6ihFANRmk7Rce8k_A8e")
	viper.SetDefault("OKTA_SINCE", "2026-02-16T00:00:00Z")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DATABASE", "okta")
	viper.SetDefault("WEBHOOK_URL", "http://localhost:8082/webhook")

	var cfg Config

	cfg.Okta.Domain = viper.GetString("OKTA_DOMAIN")
	cfg.Okta.Token = viper.GetString("OKTA_TOKEN")
	cfg.Okta.Since = viper.GetString("OKTA_SINCE")
	cfg.Okta.Limit = viper.GetInt("OKTA_LIMIT")

	cfg.PollInterval, _ = time.ParseDuration(viper.GetString("POLL_INTERVAL"))

	cfg.Mongo.URI = viper.GetString("MONGO_URI")
	cfg.Mongo.Database = viper.GetString("MONGO_DATABASE")
	cfg.Mongo.LogsCol = viper.GetString("MONGO_LOGS_COL")
	cfg.Mongo.StateCol = viper.GetString("MONGO_STATE_COL")

	cfg.Webhook.URL = viper.GetString("WEBHOOK_URL")
	cfg.Webhook.Method = viper.GetString("WEBHOOK_METHOD")
	cfg.Webhook.Timeout, _ = time.ParseDuration(viper.GetString("WEBHOOK_TIMEOUT"))
	cfg.Webhook.Retries = viper.GetInt("WEBHOOK_RETRIES")

	return &cfg, nil
}
