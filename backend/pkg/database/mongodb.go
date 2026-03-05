package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

// Migrate ensures all required indexes exist. Safe to call on every startup —
// MongoDB's CreateOne is idempotent and will not recreate an existing index.
func (m *MongoDB) Migrate(ctx context.Context) error {
	type indexSpec struct {
		collection string
		name       string
		model      mongo.IndexModel
	}

	specs := []indexSpec{
		// ── users ──────────────────────────────────────────────────────────────
		{
			collection: "users",
			name:       "users_google_id_unique",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "google_id", Value: 1}},
				Options: options.Index().SetUnique(true).SetSparse(true).SetName("users_google_id_unique"),
			},
		},
		{
			collection: "users",
			name:       "users_email_unique",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("users_email_unique"),
			},
		},

		// ── cinemas ────────────────────────────────────────────────────────────
		// Used by list-cinemas date-range filter and sort
		{
			collection: "cinemas",
			name:       "cinemas_start_time",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "start_time", Value: 1}},
				Options: options.Index().SetName("cinemas_start_time"),
			},
		},

		// ── bookings ───────────────────────────────────────────────────────────
		{
			collection: "bookings",
			name:       "bookings_user_id",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "user_id", Value: 1}},
				Options: options.Index().SetName("bookings_user_id"),
			},
		},
		{
			collection: "bookings",
			name:       "bookings_cinema_id",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "cinema_id", Value: 1}},
				Options: options.Index().SetName("bookings_cinema_id"),
			},
		},
		// Used by the timeout ticker to find PENDING bookings efficiently
		{
			collection: "bookings",
			name:       "bookings_status_created_at",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: 1}},
				Options: options.Index().SetName("bookings_status_created_at"),
			},
		},

		// ── audit_logs ─────────────────────────────────────────────────────────
		{
			collection: "audit_logs",
			name:       "audit_logs_timestamp_desc",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "timestamp", Value: -1}},
				Options: options.Index().SetName("audit_logs_timestamp_desc"),
			},
		},
		{
			collection: "audit_logs",
			name:       "audit_logs_user_id",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "user_id", Value: 1}},
				Options: options.Index().SetName("audit_logs_user_id"),
			},
		},
		{
			collection: "audit_logs",
			name:       "audit_logs_cinema_id",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "cinema_id", Value: 1}},
				Options: options.Index().SetName("audit_logs_cinema_id"),
			},
		},
	}

	for _, s := range specs {
		if _, err := m.DB.Collection(s.collection).Indexes().CreateOne(ctx, s.model); err != nil {
			return fmt.Errorf("create index on %s: %w", s.collection, err)
		}
		log.Printf("  index ensured: %s / %s", s.collection, s.name)
	}

	log.Println("MongoDB migration completed")
	return nil
}
