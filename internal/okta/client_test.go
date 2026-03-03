package okta

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseNextLink(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		expected    string
	}{
		{
			name:        "Valid next link",
			headerValue: `<https://trial-123.okta.com/api/v1/logs?after=123>; rel="next"`,
			expected:    "https://trial-123.okta.com/api/v1/logs?after=123",
		},
		{
			name:        "Missing next link",
			headerValue: `<https://trial-123.okta.com/api/v1/logs?after=abc>; rel="self"`,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseNextLink(tt.headerValue)
			if result != tt.expected {
				t.Errorf("parseNextLink() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFetchLogs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Link", `<https://mock.okta.com/next>; rel="next"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"uuid": "12345", "published": "2026-03-02T12:00:00.000Z"}]`))
	}))
	defer server.Close()

	client := NewClient("mock-token")
	resp, err := client.FetchLogs(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resp.Logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(resp.Logs))
	}
	if resp.NextURL != "https://mock.okta.com/next" {
		t.Errorf("Expected next URL, got %s", resp.NextURL)
	}
}
