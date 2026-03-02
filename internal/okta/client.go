package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"okta-ingestor/internal/models"
)

type Client struct {
	http   *http.Client
	domain string
	token  string
	limit  int
}

func New(domain, token string, limit int) *Client {
	return &Client{
		http:   &http.Client{Timeout: 30 * time.Second},
		domain: domain,
		token:  token,
		limit:  limit,
	}
}

func (c *Client) FetchLogs(ctx context.Context, since time.Time) ([]models.LogEvent, time.Time, error) {
	var all []models.LogEvent
	maxPublished := since

	url := fmt.Sprintf("https://%s/api/v1/logs?since=%s&limit=%d&sortOrder=ASCENDING",
		c.domain,
		since.Format(time.RFC3339Nano),
		c.limit)

	// slog.Info("fetching logs", "url", url)

	for {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("Authorization", "SSWS "+c.token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.http.Do(req)
		slog.Info("fetching logs", "url", resp.Request.URL.String(), "status", resp.StatusCode)
		if err != nil {
			return nil, since, err
		}

		defer resp.Body.Close()

		if resp.StatusCode == 429 {
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			return nil, since, fmt.Errorf("okta error: %s", resp.Status)
		}

		var raw []map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return nil, since, err
		}

		if len(raw) == 0 {
			break
		}

		links := strings.Split(resp.Header.Get("Link"), ",")
		next := ""
		for _, l := range links {
			if strings.Contains(l, `rel="next"`) {
				// Extract the URL between < and >
				start := strings.Index(l, "<")
				end := strings.Index(l, ">")
				if start != -1 && end != -1 {
					next = l[start+1 : end]
					break
				}
			}
		}
		if strings.Contains(next, "after=") {
			slog.Info("pagination detected, but skipping due to time-based polling")
			return nil, since, fmt.Errorf("completed")
		}

		for _, r := range raw {
			pubStr := r["published"].(string)
			pub, _ := time.Parse(time.RFC3339Nano, pubStr)

			ev := models.LogEvent{
				UUID:      r["uuid"].(string),
				Published: pub,
				EventType: fmt.Sprint(r["eventType"]),
				Severity:  fmt.Sprint(r["severity"]),
				Raw:       r,
			}

			if pub.After(maxPublished) {
				maxPublished = pub
			}

			all = append(all, ev)
		}

		nextURL := extractNextLink(resp.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		url = nextURL
	}

	return all, maxPublished, nil
}

func extractNextLink(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	links := strings.Split(linkHeader, ",")
	for _, l := range links {
		if strings.Contains(l, `rel="next"`) {
			// Extract the URL between < and >
			start := strings.Index(l, "<")
			end := strings.Index(l, ">")
			if start != -1 && end != -1 {
				return l[start+1 : end]
			}
		}
	}
	return ""
}
