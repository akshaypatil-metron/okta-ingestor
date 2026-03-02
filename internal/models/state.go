package models

// SyncState represents the last fetched cursor URL to prevent data loss or duplication.
type SyncState struct {
	ID      string `bson:"_id"`
	NextURL string `bson:"next_url"`
}
