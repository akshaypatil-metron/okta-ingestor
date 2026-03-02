package models

import "time"

type State struct {
	ID        string    `bson:"_id"`
	LastSince time.Time `bson:"last_since"`
}
