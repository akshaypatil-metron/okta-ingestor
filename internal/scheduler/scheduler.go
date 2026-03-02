package scheduler

import (
	"log"
	"time"

	"okta-ingestor/internal/processor"
)

type Scheduler struct {
	processor *processor.Processor
	interval  time.Duration
}

func NewScheduler(proc *processor.Processor, interval time.Duration) *Scheduler {
	return &Scheduler{
		processor: proc,
		interval:  interval,
	}
}

func (s *Scheduler) Start() {
	log.Printf("Starting scheduler. Polling every %v...", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Execute immediately on startup
	s.processor.ProcessBatch()

	// Wait for ticker for subsequent executions
	for range ticker.C {
		s.processor.ProcessBatch()
	}
}
