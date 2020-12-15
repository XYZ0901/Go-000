package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Uri   string
	Retry int32
	Ctx   context.Context
}

func NewConfig(ctx context.Context, uri string) Config {
	return Config{
		Uri:   uri,
		Retry: 0,
		Ctx:   ctx,
	}
}

func CreateClient(conf Config) (*mongo.Client, error) {
	conn, err := mongo.Connect(conf.Ctx, options.Client().ApplyURI(conf.Uri))
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(conf.Ctx, nil); err != nil {
		if conf.Retry >= 3 {
			return nil, err
		}
		conf.Retry++
		time.Sleep(time.Second * 2)
		return CreateClient(conf)
	}
	return conn, err
}

func GetCollection() {

}
