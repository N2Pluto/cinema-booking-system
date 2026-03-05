package repository

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AdminBookingFilter holds optional filters for the admin booking list.
type AdminBookingFilter struct {
	Movie  string // partial match on cinema.movie_name
	Date   string // "YYYY-MM-DD" — filters by created_at day
	Status string // exact BookingStatus value
	Page   int
	Limit  int
}

// BookingWithCinema is a booking enriched with cinema info via aggregation $lookup.
type BookingWithCinema struct {
	entity.Booking `bson:",inline"`
	MovieName      string    `bson:"movie_name" json:"movie_name"`
	TheaterNo      int       `bson:"theater_no" json:"theater_no"`
	StartTime      time.Time `bson:"start_time" json:"start_time"`
	EndTime        time.Time `bson:"end_time"   json:"end_time"`
}

// AdminBookingResult is a paginated list of enriched bookings.
type AdminBookingResult struct {
	Data       []*BookingWithCinema `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) (*entity.Booking, error)
	FindByID(ctx context.Context, id string) (*entity.Booking, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error)
	UpdateStatus(ctx context.Context, id bson.ObjectID, status entity.BookingStatus) error
	// FindExpiredPending returns PENDING bookings created before the given time.
	FindExpiredPending(ctx context.Context, before time.Time) ([]*entity.Booking, error)
	// FindAll returns all bookings with optional filters and cinema info (admin only).
	FindAll(ctx context.Context, f AdminBookingFilter) (*AdminBookingResult, error)
}
