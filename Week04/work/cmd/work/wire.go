// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"work/internal/biz"
	"work/internal/data"
	"work/internal/pkg/database/mongodb"
	"work/internal/service"
)

func InitProtoServe(collection *mongo.Collection) *service.Hello {
	wire.Build(data.NewMongoRepo, biz.NewHelloRequest, service.NewHello)
	return &service.Hello{}
}

func InitClient(ctx context.Context, uri string) (*mongo.Client, error) {
	wire.Build(mongodb.NewConfig, mongodb.CreateClient)
	return nil, nil
}
