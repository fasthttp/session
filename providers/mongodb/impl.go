package mongodb

import (
	"context"
	"time"

	"github.com/savsgio/gotils/strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (p *Provider) getCollection() *mongo.Collection {
	return p.db.Database(p.config.Database).Collection(p.config.Collection)
}

func (p *Provider) getSessionId(sessionID []byte) string {
	return strconv.B2S(sessionID)
}

func (p *Provider) getFilter(id []byte) primitive.D {
	return bson.D{{Key: "sessionId", Value: p.getSessionId(id)}}
}

// Get returns the session value stored in the database.
func (p *Provider) Get(id []byte) ([]byte, error) {

	var i item
	err := p.getCollection().FindOne(context.Background(), p.getFilter(id)).Decode(&i)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil // Session not found
		}
		return nil, err
	}

	return i.Data, nil
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id []byte, data []byte, expiration time.Duration) error {
	sessionId := p.getSessionId(id)

	_, err := p.getCollection().UpdateOne(context.Background(), p.getFilter(id), bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "sessionId", Value: sessionId},
			{Key: "data", Value: data},
			{Key: "expiration", Value: expiration},
		}},
	}, options.Update().SetUpsert(true))

	return err
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {

	_, err := p.getCollection().DeleteOne(context.Background(), p.getFilter(id))

	return err
}

// Regenerate updates the session id and expiration with the new session id
// of the the given current session id
func (p *Provider) Regenerate(id []byte, newID []byte, expiration time.Duration) error {
	var i item
	err := p.getCollection().FindOneAndDelete(context.Background(), p.getFilter(id)).Decode(&i)
	if err != nil {
		return err
	}

	_, err = p.getCollection().InsertOne(context.Background(), bson.D{
		{Key: "sessionId", Value: p.getSessionId(newID)},
		{Key: "data", Value: i.Data},
		{Key: "expiration", Value: expiration},
	})

	return err
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	count, err := p.getCollection().CountDocuments(context.Background(), bson.D{})

	if err != nil {
		return 0
	}

	return int(count)
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return true
}

// GC destroys the expired sessions
func (p *Provider) GC() error {
	filter := bson.D{{Key: "expiration", Value: bson.D{{Key: "$lt", Value: time.Now().Unix()}}}}

	_, err := p.getCollection().DeleteMany(context.Background(), filter)

	return err
}
