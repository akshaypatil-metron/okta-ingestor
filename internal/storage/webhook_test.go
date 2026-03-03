package storage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"okta-ingestor/internal/models"
)

func TestWebhookClient_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWebhookClient(server.URL)
	err := client.Send([]models.LogEvent{map[string]interface{}{"event": "test_event"}})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWebhookClient_Send_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewWebhookClient(server.URL)
	err := client.Send([]models.LogEvent{map[string]interface{}{"event": "test_event"}})

	if err == nil {
		t.Errorf("Expected an error for 500 status code, but got nil")
	}
}
