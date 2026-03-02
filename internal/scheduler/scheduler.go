package scheduler

import (
	"context"
	"log/slog"
	"time"

	"okta-ingestor/internal/models"
	"okta-ingestor/internal/okta"
	"okta-ingestor/internal/processor"
	"okta-ingestor/internal/storage"
)

type Scheduler struct {
	okta   *okta.Client
	repo   *storage.MongoRepository
	hook   *storage.WebhookSender
	ticker *time.Ticker
}

func New(okta *okta.Client, repo *storage.MongoRepository, hook *storage.WebhookSender, interval time.Duration) *Scheduler {
	return &Scheduler{
		okta:   okta,
		repo:   repo,
		hook:   hook,
		ticker: time.NewTicker(interval),
	}
}

func (s *Scheduler) Start(ctx context.Context, initialSince time.Time) {
	since := initialSince

	for {
		select {
		case <-ctx.Done():
			slog.Info("scheduler stopped")
			return

		case <-s.ticker.C:
			slog.Info("polling okta", "since", since)

			events, maxPub, err := s.okta.FetchLogs(ctx, since)
			if err != nil {
				slog.Error("fetch failed", "error", err)
				continue
			}

			events = processor.Deduplicate(events)

			if err := s.repo.InsertLogs(ctx, events); err != nil {
				slog.Error("mongo insert failed", "error", err)
				continue
			}

			if err := s.hook.Send(ctx, events); err != nil {
				slog.Error("webhook failed", "error", err)
				continue
			}

			since = maxPub
			s.repo.SaveState(ctx, models.State{LastSince: maxPub})
		}
	}
}
