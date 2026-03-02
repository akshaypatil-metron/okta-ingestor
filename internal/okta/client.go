package okta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"okta-ingestor/internal/models"
)

type Client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) FetchLogs(url string) (*models.OktaResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "SSWS "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("okta API returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a map so we can access the Okta uuid directly
	var rawLogs []map[string]interface{}
	if err := json.Unmarshal(body, &rawLogs); err != nil {
		return nil, err
	}

	var logs []models.LogEvent
	for _, entry := range rawLogs {
		// THIS IS THE MAGIC TRICK:
		// We tell MongoDB to use Okta's 'uuid' as the main document _id.
		// This mathematically guarantees 0 duplicates in your database.
		if uuid, ok := entry["uuid"].(string); ok {
			entry["_id"] = uuid
		}
		logs = append(logs, entry)
	}

	// Look through ALL headers Okta sent to find the "next" link
	var nextURL string
	for _, linkHeader := range resp.Header.Values("Link") {
		parsed := parseNextLink(linkHeader)
		if parsed != "" {
			nextURL = parsed
			break
		}
	}

	return &models.OktaResponse{Logs: logs, NextURL: nextURL}, nil
}

func parseNextLink(headerValue string) string {
	links := strings.Split(headerValue, ",")
	for _, link := range links {
		if strings.Contains(link, `rel="next"`) || strings.Contains(link, `rel=next`) {
			parts := strings.Split(link, ";")
			if len(parts) > 0 {
				return strings.Trim(strings.TrimSpace(parts[0]), "<>")
			}
		}
	}
	return ""
}
