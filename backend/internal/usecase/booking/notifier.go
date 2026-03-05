package booking

import (
	"context"
	"log"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

// LoggerNotificationService is a simple mock implementation that just logs events.
// It satisfies the NotificationService interface and can be replaced with
// real Email / Line providers in the future.
type LoggerNotificationService struct{}

func (LoggerNotificationService) BookingSuccess(ctx context.Context, b *entity.Booking) error {
	log.Printf("[notification] booking success: booking_id=%s user_id=%s cinema_id=%s seats=%v",
		b.ID.Hex(), b.UserID.Hex(), b.CinemaID.Hex(), b.SeatNumbers)
	return nil
}

func (LoggerNotificationService) BookingTimeout(ctx context.Context, b *entity.Booking) error {
	log.Printf("[notification] booking timeout: booking_id=%s user_id=%s cinema_id=%s seats=%v",
		b.ID.Hex(), b.UserID.Hex(), b.CinemaID.Hex(), b.SeatNumbers)
	return nil
}

