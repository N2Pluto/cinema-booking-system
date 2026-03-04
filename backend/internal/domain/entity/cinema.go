package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatLocked    SeatStatus = "LOCKED"
	SeatBooked    SeatStatus = "BOOKED"
)

// Seat is embedded inside Cinema document.
type Seat struct {
	SeatNo string        `bson:"seat_no" json:"seat_no"`
	Status SeatStatus    `bson:"status"  json:"status"`
	Price  int           `bson:"price"   json:"price"`
	UserID *bson.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type Cinema struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieName string        `bson:"movie_name"    json:"movie_name"`
	TheaterNo int           `bson:"theater_no"    json:"theater_no"`
	StartTime time.Time     `bson:"start_time"    json:"start_time"`
	EndTime   time.Time     `bson:"end_time"      json:"end_time"`
	Price     int           `bson:"price"         json:"price"`
	Seats     []Seat        `bson:"seats"         json:"seats"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
}
