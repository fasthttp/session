package mongodb

import "fmt"

func newErrMongoConnection(err error) error {
	return fmt.Errorf("MongoDB connection error: %w", err)
}
