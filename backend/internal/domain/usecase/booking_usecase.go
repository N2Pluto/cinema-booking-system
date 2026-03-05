package usecase

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
)

type LockSeatsInput struct {
	UserID      string
	CinemaID    string
	SeatNumbers []string
}

type AdminListBookingsInput struct {
	Movie  string
	Date   string // "YYYY-MM-DD"
	Status string
	Page   int
	Limit  int
}

type BookingUseCase interface {
	LockSeats(ctx context.Context, in LockSeatsInput) (*entity.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID string, userID string) (*entity.Booking, error)
	CancelBooking(ctx context.Context, bookingID string, userID string) error
	GetMyBookings(ctx context.Context, userID string) ([]*entity.Booking, error)
	ListBookings(ctx context.Context, in AdminListBookingsInput) (*repository.AdminBookingResult, error)
}
