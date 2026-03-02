package models

import "time"

type LogEvent struct {
	UUID      string    `bson:"uuid" json:"uuid"`
	Published time.Time `bson:"published" json:"published"`
	EventType string    `bson:"eventType" json:"eventType"`
	Severity  string    `bson:"severity" json:"severity"`
	Raw       any       `bson:"raw" json:"raw"`
}
