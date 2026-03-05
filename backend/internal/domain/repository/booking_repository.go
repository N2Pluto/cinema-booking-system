package repository

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) (*entity.Booking, error)
	FindByID(ctx context.Context, id string) (*entity.Booking, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error)
	UpdateStatus(ctx context.Context, id bson.ObjectID, status entity.BookingStatus) error
	// FindExpiredPending returns PENDING bookings created before the given time.
	FindExpiredPending(ctx context.Context, before time.Time) ([]*entity.Booking, error)
}
