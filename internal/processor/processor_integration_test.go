//go:build integration

package processor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"okta-ingestor/internal/okta"
	"okta-ingestor/internal/storage"
)

func TestProcessorIntegration(t *testing.T) {
	// Setup shared MongoDB connection for all tests
	mongoURI := "mongodb://localhost:27017"
	mongoStore, err := storage.NewMongoStore(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to Mongo: %v", err)
	}

	// Fake Webhook (always succeeds)
	webhookHit := false
	webhookServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookHit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer webhookServer.Close()
	webhookClient := storage.NewWebhookClient(webhookServer.URL)

	t.Run("Happy Path - Logs Found", func(t *testing.T) {
		webhookHit = false // reset
		oktaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Link", `<http://fake-next.com>; rel="next"`)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"uuid": "processor-test-uuid-1", "published": "2026-03-03T10:00:00.000Z"}]`))
		}))
		defer oktaServer.Close()

		mongoStore.SaveCursor(oktaServer.URL) // Reset cursor
		client := okta.NewClient("dummy-token")
		proc := NewProcessor(client, mongoStore, webhookClient, oktaServer.URL, 100, "")

		proc.ProcessBatch()

		if !webhookHit {
			t.Errorf("Expected webhook to be hit, but it wasn't")
		}
	})

	t.Run("Edge Case - No Logs Found", func(t *testing.T) {
		oktaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`)) // Return empty array
		}))
		defer oktaServer.Close()

		mongoStore.SaveCursor(oktaServer.URL) // Reset cursor
		client := okta.NewClient("dummy-token")
		proc := NewProcessor(client, mongoStore, webhookClient, oktaServer.URL, 100, "")

		proc.ProcessBatch()
		// This will successfully hit your red `else` block: "No new logs found!"
	})

	t.Run("Edge Case - Okta API Error", func(t *testing.T) {
		oktaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized) // Simulate a 401 Error
		}))
		defer oktaServer.Close()

		mongoStore.SaveCursor(oktaServer.URL) // Reset cursor
		client := okta.NewClient("dummy-token")
		proc := NewProcessor(client, mongoStore, webhookClient, oktaServer.URL, 100, "")

		proc.ProcessBatch()
		// This will successfully hit your red `if err != nil` block at the top!
	})
}
