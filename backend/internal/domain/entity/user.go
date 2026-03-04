package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleID    string        `bson:"google_id"     json:"google_id"`
	Email       string        `bson:"email"         json:"email"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Role        UserRole      `bson:"role"          json:"role"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
}
