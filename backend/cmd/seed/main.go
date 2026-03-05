package main

import (
	"context"
	"fmt"
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

// makeSeats generates rows×cols seats.
// bookedNos and lockedNos are seat numbers pre-set to those statuses for demo purposes.
func makeSeats(rows int, cols int, basePrice int, premiumRows int, bookedNos, lockedNos map[string]bool) []entity.Seat {
	seats := make([]entity.Seat, 0, rows*cols)
	for r := 0; r < rows; r++ {
		rowLetter := string(rune('A' + r))
		price := basePrice
		if r >= rows-premiumRows { // last N rows are premium
			price = basePrice + 50
		}
		for c := 1; c <= cols; c++ {
			seatNo := rowLetter + fmt.Sprintf("%d", c)
			status := entity.SeatAvailable
			if bookedNos[seatNo] {
				status = entity.SeatBooked
			} else if lockedNos[seatNo] {
				status = entity.SeatLocked
			}
			seats = append(seats, entity.Seat{SeatNo: seatNo, Status: status, Price: price})
		}
	}
	return seats
}

func seedCinemas(ctx context.Context, db *database.MongoDB) {
	col := db.Collection("cinemas")
	col.Drop(ctx)

	now := time.Now()

	// Helper: set of seat numbers
	s := func(nos ...string) map[string]bool {
		m := make(map[string]bool, len(nos))
		for _, n := range nos {
			m[n] = true
		}
		return m
	}

	cinemas := []interface{}{
		// ── โรง 1 ────────────────────────────────────────────────────────────
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Avengers: Endgame",
			TheaterNo: 1,
			StartTime: now.Add(1 * time.Hour),
			EndTime:   now.Add(4 * time.Hour),
			Price:     200,
			Seats:     makeSeats(5, 10, 200, 1, s("A1", "A2", "A5", "B3", "B4", "C1"), s("D2", "D3")),
			CreatedAt: now,
		},
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Avengers: Endgame",
			TheaterNo: 1,
			StartTime: now.Add(6 * time.Hour),
			EndTime:   now.Add(9 * time.Hour),
			Price:     200,
			Seats:     makeSeats(5, 10, 200, 1, s("B1", "B2"), s()),
			CreatedAt: now,
		},

		// ── โรง 2 ────────────────────────────────────────────────────────────
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Interstellar",
			TheaterNo: 2,
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(5*time.Hour + 49*time.Minute),
			Price:     220,
			Seats:     makeSeats(6, 8, 220, 2, s("A1", "A2", "A3", "B5", "B6", "C2", "C3", "C4"), s("E1", "E2")),
			CreatedAt: now,
		},
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Dune: Part Two",
			TheaterNo: 2,
			StartTime: now.Add(7 * time.Hour),
			EndTime:   now.Add(9*time.Hour + 46*time.Minute),
			Price:     250,
			Seats:     makeSeats(6, 8, 250, 2, s("A1", "A2", "B1", "B2"), s("C5", "C6")),
			CreatedAt: now,
		},

		// ── โรง 3 ────────────────────────────────────────────────────────────
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Oppenheimer",
			TheaterNo: 3,
			StartTime: now.Add(3 * time.Hour),
			EndTime:   now.Add(6*time.Hour + 30*time.Minute),
			Price:     180,
			Seats:     makeSeats(4, 12, 180, 1, s("A1", "A2", "A3", "A4", "B6", "B7"), s()),
			CreatedAt: now,
		},
		entity.Cinema{
			ID:        bson.NewObjectID(),
			MovieName: "Spider-Man: No Way Home",
			TheaterNo: 3,
			StartTime: now.Add(8 * time.Hour),
			EndTime:   now.Add(10*time.Hour + 28*time.Minute),
			Price:     180,
			Seats:     makeSeats(4, 12, 180, 1, s(), s()),
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
