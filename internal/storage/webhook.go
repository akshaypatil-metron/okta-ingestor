package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"okta-ingestor/internal/models"
)

type WebhookSender struct {
	client  *http.Client
	url     string
	method  string
	retries int
}

func NewWebhookSender(url, method string, timeout time.Duration, retries int) *WebhookSender {
	return &WebhookSender{
		client:  &http.Client{Timeout: timeout},
		url:     url,
		method:  method,
		retries: retries,
	}
}

func (w *WebhookSender) Send(ctx context.Context, logs []models.LogEvent) error {
	body, _ := json.Marshal(logs)

	for i := 0; i < w.retries; i++ {
		req, _ := http.NewRequestWithContext(ctx, w.method, w.url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := w.client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(time.Second)
	}

	return nil
}