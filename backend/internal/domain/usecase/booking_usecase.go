package usecase

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type LockSeatsInput struct {
	UserID      string
	CinemaID    string
	SeatNumbers []string
}

type BookingUseCase interface {
	LockSeats(ctx context.Context, in LockSeatsInput) (*entity.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID string, userID string) (*entity.Booking, error)
	GetMyBookings(ctx context.Context, userID string) ([]*entity.Booking, error)
}
