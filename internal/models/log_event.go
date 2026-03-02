package models

// LogEvent represents a single Okta log entry.
type LogEvent interface{}

// OktaResponse holds the parsed logs and the pagination cursor.
type OktaResponse struct {
	Logs    []LogEvent
	NextURL string
}
