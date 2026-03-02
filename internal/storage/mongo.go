package storage

import (
	"context"
	"errors"
	"okta-ingestor/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	logs  *mongo.Collection
	state *mongo.Collection
}

func NewMongoRepository(db *mongo.Database, logsCol, stateCol string) *MongoRepository {
	return &MongoRepository{
		logs:  db.Collection(logsCol),
		state: db.Collection(stateCol),
	}
}

func (r *MongoRepository) InsertLogs(ctx context.Context, logs []models.LogEvent) error {
	docs := make([]any, len(logs))
	for i := range logs {
		docs[i] = logs[i]
	}

	_, err := r.logs.InsertMany(ctx, docs)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return nil
	}
	return err
}

func (r *MongoRepository) LoadState(ctx context.Context) (models.State, error) {
	var s models.State
	err := r.state.FindOne(ctx, bson.M{"_id": "cursor"}).Decode(&s)
	return s, err
}

func (r *MongoRepository) SaveState(ctx context.Context, since models.State) error {
	_, err := r.state.UpdateByID(ctx, "cursor",
		bson.M{"$set": bson.M{"last_since": since.LastSince}},
	)
	return err
}

var ErrNoState = errors.New("no state found")
