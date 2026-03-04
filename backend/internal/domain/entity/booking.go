package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "PENDING"
	BookingSuccess   BookingStatus = "SUCCESS"
	BookingTimeout   BookingStatus = "TIMEOUT"
	BookingCancelled BookingStatus = "CANCELLED"
)

type Booking struct {
	ID          bson.ObjectID `bson:"_id,omitempty"  json:"id"`
	UserID      bson.ObjectID `bson:"user_id"        json:"user_id"`
	CinemaID    bson.ObjectID `bson:"cinema_id"      json:"cinema_id"`
	SeatNumbers []string      `bson:"seat_numbers"   json:"seat_numbers"`
	TotalAmount float64       `bson:"total_amount"   json:"total_amount"`
	Status      BookingStatus `bson:"status"         json:"status"`
	CreatedAt   time.Time     `bson:"created_at"     json:"created_at"`
}
