package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"okta-ingestor/internal/models"
)

type WebhookClient struct {
	URL        string
	HTTPClient *http.Client
}

func NewWebhookClient(url string) *WebhookClient {
	return &WebhookClient{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (w *WebhookClient) Send(logs []models.LogEvent) error {
	jsonData, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status: %s", resp.Status)
	}
	return nil
}
