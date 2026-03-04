package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/ws"
	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const lockTTL = 5 * time.Minute

var _ domainusecase.BookingUseCase = (*UseCase)(nil)

type UseCase struct {
	cinemaRepo  repository.CinemaRepository
	bookingRepo repository.BookingRepository
	auditRepo   repository.AuditLogRepository
	hub         *ws.Hub
}

func NewUseCase(
	cinemaRepo repository.CinemaRepository,
	bookingRepo repository.BookingRepository,
	auditRepo repository.AuditLogRepository,
	hub *ws.Hub,
) *UseCase {
	return &UseCase{
		cinemaRepo:  cinemaRepo,
		bookingRepo: bookingRepo,
		auditRepo:   auditRepo,
		hub:         hub,
	}
}

// StartTimeoutTicker runs a background goroutine that releases expired PENDING bookings every 30s.
func (uc *UseCase) StartTimeoutTicker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				uc.releaseExpired(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (uc *UseCase) releaseExpired(ctx context.Context) {
	before := time.Now().Add(-lockTTL)
	bookings, err := uc.bookingRepo.FindExpiredPending(ctx, before)
	if err != nil {
		log.Printf("releaseExpired find: %v", err)
		return
	}

	for _, b := range bookings {
		if err := uc.cinemaRepo.ReleaseSeats(ctx, b.CinemaID, b.SeatNumbers); err != nil {
			log.Printf("releaseExpired releaseSeats %s: %v", b.ID.Hex(), err)
			continue
		}
		if err := uc.bookingRepo.UpdateStatus(ctx, b.ID, entity.BookingTimeout); err != nil {
			log.Printf("releaseExpired updateStatus %s: %v", b.ID.Hex(), err)
		}
		_ = uc.auditRepo.Create(ctx, &entity.AuditLog{
			EventType: entity.EventBookingTimeout,
			UserID:    b.UserID,
			CinemaID:  b.CinemaID,
			Details:   fmt.Sprintf("seats %v released after timeout", b.SeatNumbers),
		})
		uc.broadcastSeats(ctx, b.CinemaID.Hex())
	}
}

func (uc *UseCase) LockSeats(ctx context.Context, in domainusecase.LockSeatsInput) (*entity.Booking, error) {
	userOID, err := bson.ObjectIDFromHex(in.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	cinemaOID, err := bson.ObjectIDFromHex(in.CinemaID)
	if err != nil {
		return nil, fmt.Errorf("invalid cinema id: %w", err)
	}
	if len(in.SeatNumbers) == 0 {
		return nil, fmt.Errorf("no seats selected")
	}

	cinema, err := uc.cinemaRepo.FindByID(ctx, in.CinemaID)
	if err != nil || cinema == nil {
		return nil, fmt.Errorf("cinema not found")
	}

	// Calculate total amount from seat prices
	seatPriceMap := make(map[string]int, len(cinema.Seats))
	for _, s := range cinema.Seats {
		seatPriceMap[s.SeatNo] = s.Price
	}
	var total float64
	for _, sno := range in.SeatNumbers {
		p, ok := seatPriceMap[sno]
		if !ok {
			return nil, fmt.Errorf("seat %s not found in cinema", sno)
		}
		total += float64(p)
	}

	// Atomically lock seats
	if err := uc.cinemaRepo.LockSeats(ctx, cinemaOID, in.SeatNumbers, userOID); err != nil {
		return nil, fmt.Errorf("lock seats: %w", err)
	}

	booking, err := uc.bookingRepo.Create(ctx, &entity.Booking{
		UserID:      userOID,
		CinemaID:    cinemaOID,
		SeatNumbers: in.SeatNumbers,
		TotalAmount: total,
		Status:      entity.BookingPending,
	})
	if err != nil {
		// Best-effort rollback
		_ = uc.cinemaRepo.ReleaseSeats(ctx, cinemaOID, in.SeatNumbers)
		return nil, fmt.Errorf("create booking: %w", err)
	}

	_ = uc.auditRepo.Create(ctx, &entity.AuditLog{
		EventType: entity.EventSeatLocked,
		UserID:    userOID,
		CinemaID:  cinemaOID,
		Details:   fmt.Sprintf("seats %v locked by user %s", in.SeatNumbers, in.UserID),
	})

	uc.broadcastSeats(ctx, in.CinemaID)
	return booking, nil
}

func (uc *UseCase) ConfirmBooking(ctx context.Context, bookingID string, userID string) (*entity.Booking, error) {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil || booking == nil {
		return nil, fmt.Errorf("booking not found")
	}
	if booking.UserID.Hex() != userID {
		return nil, fmt.Errorf("booking does not belong to user")
	}
	if booking.Status != entity.BookingPending {
		return nil, fmt.Errorf("booking is not in PENDING state")
	}

	if err := uc.cinemaRepo.ConfirmSeats(ctx, booking.CinemaID, booking.SeatNumbers); err != nil {
		return nil, fmt.Errorf("confirm seats: %w", err)
	}
	if err := uc.bookingRepo.UpdateStatus(ctx, booking.ID, entity.BookingSuccess); err != nil {
		return nil, fmt.Errorf("update booking: %w", err)
	}
	booking.Status = entity.BookingSuccess

	_ = uc.auditRepo.Create(ctx, &entity.AuditLog{
		EventType: entity.EventBookingSuccess,
		UserID:    booking.UserID,
		CinemaID:  booking.CinemaID,
		Details:   fmt.Sprintf("booking %s confirmed, seats %v", bookingID, booking.SeatNumbers),
	})

	uc.broadcastSeats(ctx, booking.CinemaID.Hex())
	return booking, nil
}

func (uc *UseCase) GetMyBookings(ctx context.Context, userID string) ([]*entity.Booking, error) {
	return uc.bookingRepo.FindByUserID(ctx, userID)
}

// broadcastSeats fetches the latest cinema and pushes its seat list to all WS clients.
func (uc *UseCase) broadcastSeats(ctx context.Context, cinemaID string) {
	cinema, err := uc.cinemaRepo.FindByID(ctx, cinemaID)
	if err != nil || cinema == nil {
		return
	}
	payload := map[string]interface{}{
		"type":  "seat_update",
		"seats": cinema.Seats,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	uc.hub.Broadcast <- &ws.BroadcastMsg{
		CinemaID: cinemaID,
		Payload:  data,
	}
}
