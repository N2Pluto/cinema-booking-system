package booking

import (
	"context"
	"testing"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// simple in-memory fakes for a happy-path LockSeats + ConfirmBooking flow

type fakeCinemaRepo struct {
	cinema *entity.Cinema
}

func (f *fakeCinemaRepo) List(ctx context.Context, _ repository.CinemaFilter) ([]*entity.Cinema, int64, error) {
	if f.cinema == nil {
		return []*entity.Cinema{}, 0, nil
	}
	return []*entity.Cinema{f.cinema}, 1, nil
}

func (f *fakeCinemaRepo) FindByID(ctx context.Context, id string) (*entity.Cinema, error) {
	return f.cinema, nil
}

func (f *fakeCinemaRepo) Create(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error) {
	f.cinema = cinema
	return cinema, nil
}

func (f *fakeCinemaRepo) LockSeats(ctx context.Context, cinemaID bson.ObjectID, seatNumbers []string, userID bson.ObjectID) error {
	for i, s := range f.cinema.Seats {
		for _, n := range seatNumbers {
			if s.SeatNo == n {
				f.cinema.Seats[i].Status = entity.SeatLocked
			}
		}
	}
	return nil
}

func (f *fakeCinemaRepo) ReleaseSeats(ctx context.Context, cinemaID bson.ObjectID, seatNumbers []string) error {
	return nil
}

func (f *fakeCinemaRepo) ConfirmSeats(ctx context.Context, cinemaID bson.ObjectID, seatNumbers []string) error {
	for i, s := range f.cinema.Seats {
		for _, n := range seatNumbers {
			if s.SeatNo == n {
				f.cinema.Seats[i].Status = entity.SeatBooked
			}
		}
	}
	return nil
}

type fakeBookingRepo struct {
	current *entity.Booking
}

func (f *fakeBookingRepo) Create(ctx context.Context, b *entity.Booking) (*entity.Booking, error) {
	b.ID = bson.NewObjectID()
	f.current = b
	return b, nil
}

func (f *fakeBookingRepo) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	return f.current, nil
}

func (f *fakeBookingRepo) UpdateStatus(ctx context.Context, id bson.ObjectID, status entity.BookingStatus) error {
	if f.current != nil {
		f.current.Status = status
	}
	return nil
}

func (f *fakeBookingRepo) FindExpiredPending(ctx context.Context, before time.Time) ([]*entity.Booking, error) {
	return nil, nil
}

func (f *fakeBookingRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	return nil, nil
}

func (f *fakeBookingRepo) FindAll(ctx context.Context, filter repository.AdminBookingFilter) (*repository.AdminBookingResult, error) {
	return nil, nil
}

type fakeSeatLocker struct{}

func (f *fakeSeatLocker) AcquireSeatLock(ctx context.Context, cinemaID, seatNo, token string, ttl time.Duration) error {
	return nil
}

func (f *fakeSeatLocker) ReleaseSeatLock(ctx context.Context, cinemaID, seatNo, token string) error {
	return nil
}

type fakePublisher struct{}

func (f *fakePublisher) PublishAuditLog(ctx context.Context, log *entity.AuditLog) error {
	return nil
}

type noopHub struct{}

func TestLockAndConfirmBooking_HappyPath(t *testing.T) {
	userID := bson.NewObjectID()
	cinemaID := bson.NewObjectID()

	cinema := &entity.Cinema{
		ID: cinemaID,
		Seats: []entity.Seat{
			{SeatNo: "A1", Price: 100, Status: entity.SeatAvailable},
			{SeatNo: "A2", Price: 100, Status: entity.SeatAvailable},
		},
	}

	uc := &UseCase{
		cinemaRepo:  &fakeCinemaRepo{cinema: cinema},
		bookingRepo: &fakeBookingRepo{},
		seatLocker:  &fakeSeatLocker{},
		publisher:   &fakePublisher{},
		hub:         nil,
	}

	ctx := context.Background()

	// Lock two seats
	b, err := uc.LockSeats(ctx, domainusecase.LockSeatsInput{
		UserID:      userID.Hex(),
		CinemaID:    cinemaID.Hex(),
		SeatNumbers: []string{"A1", "A2"},
	})
	if err != nil {
		t.Fatalf("LockSeats returned error: %v", err)
	}
	if b == nil {
		t.Fatalf("LockSeats returned nil booking")
	}
	if b.TotalAmount != 200 {
		t.Fatalf("expected total amount 200, got %v", b.TotalAmount)
	}
	if b.Status != entity.BookingPending {
		t.Fatalf("expected booking status PENDING, got %v", b.Status)
	}

	// Confirm booking
	confirmed, err := uc.ConfirmBooking(ctx, b.ID.Hex(), userID.Hex())
	if err != nil {
		t.Fatalf("ConfirmBooking returned error: %v", err)
	}
	if confirmed.Status != entity.BookingSuccess {
		t.Fatalf("expected booking status SUCCESS, got %v", confirmed.Status)
	}
}

