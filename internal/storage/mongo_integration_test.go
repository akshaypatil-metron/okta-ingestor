//go:build integration

package storage

import (
	"testing"
	"time"

	"okta-ingestor/internal/models"
)

func TestMongoStoreIntegration(t *testing.T) {
	uri := "mongodb://localhost:27017"
	store, err := NewMongoStore(uri)
	if err != nil {
		t.Fatalf("Failed to connect to Mongo: %v", err)
	}

	t.Run("Save and Get Cursor", func(t *testing.T) {
		testURL := "https://trial-5020092.okta.com/api/v1/logs?after=integration_test_123"

		if err := store.SaveCursor(testURL); err != nil {
			t.Fatalf("Failed to save cursor: %v", err)
		}

		retrievedURL := store.GetLastCursor("default_url")
		if retrievedURL != testURL {
			t.Errorf("Expected cursor %s, got %s", testURL, retrievedURL)
		}
	})

	t.Run("SaveLogs and Ignore Duplicates", func(t *testing.T) {
		log1 := map[string]interface{}{
			"_id":       "integration-uuid-999",
			"eventType": "system.test.integration",
			"published": time.Now().UTC().Format(time.RFC3339),
		}

		logs := []models.LogEvent{log1}

		// First insert should succeed
		if err := store.SaveLogs(logs); err != nil {
			t.Errorf("Expected first insert to succeed, got: %v", err)
		}

		// Second insert of exact same _id should be silently ignored
		if err := store.SaveLogs(logs); err != nil {
			t.Errorf("Expected duplicate insert to be silently ignored, got error: %v", err)
		}
	})
}
