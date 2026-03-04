package main

import (
	"context"
	"log"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/pkg/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	db, err := database.NewMongoDB()
	if err != nil {
		log.Fatalf("MongoDB connect error: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	seedUsers(ctx, db, err)
	seedCinemas(ctx, db)
	seedBookings(ctx, db)

	log.Println("Seed completed!")
}

func seedUsers(ctx context.Context, db *database.MongoDB, err error) {
	col := db.Collection("users")
	col.Drop(ctx)

	users := []interface{}{
		entity.User{
			ID:          bson.NewObjectID(),
			GoogleID:    "google-001",
			Email:       "alice@example.com",
			DisplayName: "Alice",
			Role:        entity.RoleUser,
			CreatedAt:   time.Now(),
		},
		entity.User{
			ID:          bson.NewObjectID(),
			GoogleID:    "google-002",
			Email:       "bob@example.com",
			DisplayName: "Bob",
			Role:        entity.RoleUser,
			CreatedAt:   time.Now(),
		},
		entity.User{
			ID:          bson.NewObjectID(),
			GoogleID:    "google-admin",
			Email:       "admin@cinema.com",
			DisplayName: "Admin",
			Role:        entity.RoleAdmin,
			CreatedAt:   time.Now(),
		},
	}

	_, err = col.InsertMany(ctx, users)
	if err != nil {
		log.Fatalf("seed users: %v", err)
	}
	log.Println("users seeded:", len(users))
}

func seedCinemas(ctx context.Context, db *database.MongoDB) {
	col := db.Collection("cinemas")
	col.Drop(ctx)

	now := time.Now()
	cinemas := []interface{}{
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Avengers: Endgame",
			TheaterNo: 1,
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(5 * time.Hour),
			Price:     200,
			Seats: []entity.Seat{
				{SeatNo: "A1", Status: entity.SeatAvailable, Price: 200},
				{SeatNo: "A2", Status: entity.SeatAvailable, Price: 200},
				{SeatNo: "A3", Status: entity.SeatBooked, Price: 200},
				{SeatNo: "B1", Status: entity.SeatAvailable, Price: 250},
				{SeatNo: "B2", Status: entity.SeatLocked, Price: 250},
			},
			CreatedAt: now,
		},
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Interstellar",
			TheaterNo: 2,
			StartTime: now.Add(4 * time.Hour),
			EndTime:   now.Add(7 * time.Hour),
			Price:     220,
			Seats: []entity.Seat{
				{SeatNo: "A1", Status: entity.SeatAvailable, Price: 220},
				{SeatNo: "A2", Status: entity.SeatAvailable, Price: 220},
				{SeatNo: "B1", Status: entity.SeatAvailable, Price: 270},
				{SeatNo: "B2", Status: entity.SeatAvailable, Price: 270},
			},
			CreatedAt: now,
		},
	}

	_, err := col.InsertMany(ctx, cinemas)
	if err != nil {
		log.Fatalf("seed cinemas: %v", err)
	}
	log.Println("cinemas seeded:", len(cinemas))
}

func seedBookings(ctx context.Context, db *database.MongoDB) {
	col := db.Collection("bookings")
	col.Drop(ctx)

	bookings := []interface{}{
		entity.Booking{
			ID:          bson.NewObjectID(),
			UserID:      bson.NewObjectID(),
			CinemaID:    bson.NewObjectID(),
			SeatNumbers: []string{"A3"},
			TotalAmount: 200,
			Status:      entity.BookingSuccess,
			CreatedAt:   time.Now(),
		},
	}

	_, err := col.InsertMany(ctx, bookings)
	if err != nil {
		log.Fatalf("seed bookings: %v", err)
	}
	log.Println("bookings seeded:", len(bookings))
}
