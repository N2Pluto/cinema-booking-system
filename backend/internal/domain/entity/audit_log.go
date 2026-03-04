package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventType string

const (
	EventBookingSuccess EventType = "BOOKING_SUCCESS"
	EventBookingTimeout EventType = "TIMEOUT"
	EventBookingCancel  EventType = "BOOKING_CANCELLED"
	EventSeatLocked     EventType = "SEAT_LOCKED"
	EventSeatReleased   EventType = "SEAT_RELEASED"
)

type AuditLog struct {
	ID        bson.ObjectID          `bson:"_id,omitempty" json:"id"`
	EventType EventType              `bson:"event_type"    json:"event_type"`
	UserID    bson.ObjectID          `bson:"user_id"       json:"user_id"`
	CinemaID  bson.ObjectID          `bson:"cinema_id"     json:"cinema_id"`
	Details   string                 `bson:"details"       json:"details"`
	Payload   map[string]interface{} `bson:"payload"       json:"payload"`
	Timestamp time.Time              `bson:"timestamp"     json:"timestamp"`
}
