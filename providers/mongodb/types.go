package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	ConnectionUrl string

	Database string

	Collection string
}

type Provider struct {
	config Config

	db *mongo.Client
}

type item struct {
	Data       []byte        `bson:"data"`
	Expiration time.Duration `bson:"expiration"`
	SessionId  string        `bson:"sessionId"`
}
