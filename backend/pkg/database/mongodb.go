package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB() (*MongoDB, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "cinema_booking"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Connected to MongoDB: %s / %s", uri, dbName)

	return &MongoDB{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func (m *MongoDB) Disconnect() {
	if err := m.Client.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
		return
	}
	log.Println("MongoDB disconnected")
}

// Collection returns a handle to a specific collection.
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
