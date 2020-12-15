package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(collection *mongo.Collection) *MongoRepo {
	return &MongoRepo{collection: collection}
}

func (m *MongoRepo) ReSendMessage(ctx context.Context, req string) string {
	return "hi," + req
}
