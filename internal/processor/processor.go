package processor

import (
	"fmt"
	"log"
	"time"

	"okta-ingestor/internal/okta"
	"okta-ingestor/internal/storage"
)

type Processor struct {
	oktaClient    *okta.Client
	mongoStore    *storage.MongoStore
	webhookClient *storage.WebhookClient
	baseOktaURL   string
}

func NewProcessor(oktaClient *okta.Client, mongoStore *storage.MongoStore, webhookClient *storage.WebhookClient, domain string, batchSize int, startDate string) *Processor {
	baseURL := fmt.Sprintf("%s/api/v1/logs?limit=%d", domain, batchSize)
	if startDate != "" {
		baseURL = fmt.Sprintf("%s&since=%s", baseURL, startDate)
	}

	return &Processor{
		oktaClient:    oktaClient,
		mongoStore:    mongoStore,
		webhookClient: webhookClient,
		baseOktaURL:   baseURL,
	}
}

func (p *Processor) ProcessBatch() {
	currentURL := p.mongoStore.GetLastCursor(p.baseOktaURL)

	for {
		log.Printf("Asking Okta for next batch...")
		response, err := p.oktaClient.FetchLogs(currentURL)
		if err != nil {
			log.Printf("Error fetching logs: %v", err)
			break
		}

		logCount := len(response.Logs)

		if logCount > 0 {
			// Extract the timestamp of the last log in the batch to prove it is moving forward
			var latestTimestamp string
			if lastLog, ok := response.Logs[logCount-1].(map[string]interface{}); ok {
				latestTimestamp, _ = lastLog["published"].(string)
			}

			log.Printf("Successfully fetched %d logs. Latest Timestamp: %s", logCount, latestTimestamp)

			// 1. Save to MongoDB (duplicates are silently ignored via UUID)
			if err := p.mongoStore.SaveLogs(response.Logs); err != nil {
				log.Printf("Error saving to Mongo: %v", err)
			} else {
				log.Printf("Saved to MongoDB successfully.")
			}

			// 2. Send to Webhook
			if err := p.webhookClient.Send(response.Logs); err != nil {
				log.Printf("Error sending to Webhook: %v", err)
			}

		} else {
			log.Println("No new logs found. Caught up to present time!")
			break
		}

		// 3. Move forward to the next page
		if response.NextURL != "" && response.NextURL != currentURL {
			if err := p.mongoStore.SaveCursor(response.NextURL); err != nil {
				log.Printf("Failed to save state cursor: %v", err)
			}

			currentURL = response.NextURL
			time.Sleep(1 * time.Second) // Small delay to protect Okta limits
			continue
		}

		break
	}
}
