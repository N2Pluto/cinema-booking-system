package repository

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CinemaFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	OrderBy   string // "asc" | "desc"
	Page      int
	Limit     int
}

type CinemaRepository interface {
	List(ctx context.Context, f CinemaFilter) ([]*entity.Cinema, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Cinema, error)
	Create(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error)
	// LockSeats atomically changes requested seats from AVAILABLE → LOCKED.
	// Returns error if any of the requested seats is not AVAILABLE.
	LockSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string, userID bson.ObjectID) error
	// ReleaseSeats changes seats back to AVAILABLE (used on timeout).
	ReleaseSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string) error
	// ConfirmSeats changes LOCKED → BOOKED (used on booking confirmation).
	ConfirmSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string) error
}
