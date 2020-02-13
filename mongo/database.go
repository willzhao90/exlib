package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database client and db instance
type Database struct {
	Db *mongo.Database
}

// Connect establishes mongodb connection
func Connect(ctx context.Context, v Config) *Database {
	reg := GetRegistry()

	opt := options.Client().ApplyURI(v.URI).SetRegistry(reg)
	client, err := mongo.NewClient(opt)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatalf("Cannot connect to MongoDB at %v: %v", v.URI, err)
	}
	db := client.Database(v.DbName)
	return &Database{Db: db}
}

// Close underlying mongo client and db
func (db *Database) Close(ctx context.Context) {
	db.Db.Client().Disconnect(ctx)
}

// CreateCollection instance
func (db *Database) CreateCollection(name string) *mongo.Collection {
	return db.Db.Collection(name)
}
