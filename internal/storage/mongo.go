package storage

import (
	"context"
	"time"

	"okta-ingestor/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client    *mongo.Client
	logColl   *mongo.Collection
	stateColl *mongo.Collection
}

func NewMongoStore(uri string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database("okta_data")
	store := &MongoStore{
		client:    client,
		logColl:   db.Collection("system_logs"),
		stateColl: db.Collection("sync_state"),
	}

	return store, nil
}

func (m *MongoStore) SaveLogs(logs []models.LogEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.InsertMany().SetOrdered(false)
	var insertDocs []interface{}
	for _, logItem := range logs {
		insertDocs = append(insertDocs, logItem)
	}

	_, err := m.logColl.InsertMany(ctx, insertDocs, opts)
	if err != nil {
		// Silently ignore standard Duplicate Key errors (code 11000)
		if bulkWriteErr, ok := err.(mongo.BulkWriteException); ok {
			for _, we := range bulkWriteErr.WriteErrors {
				if we.Code != 11000 {
					return err
				}
			}
			return nil
		}
		return err
	}
	return nil
}

func (m *MongoStore) GetLastCursor(defaultURL string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var state models.SyncState
	err := m.stateColl.FindOne(ctx, bson.M{"_id": "okta_cursor"}).Decode(&state)
	if err != nil || state.NextURL == "" {
		return defaultURL
	}
	return state.NextURL
}

func (m *MongoStore) SaveCursor(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": bson.M{"next_url": url}}
	_, err := m.stateColl.UpdateByID(ctx, "okta_cursor", update, opts)
	return err
}
