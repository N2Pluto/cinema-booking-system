package redis

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	"github.com/n2pluto/cinema-booking-system/pkg/redislock"
)

type seatLocker struct {
	locker *redislock.Locker
}

// NewSeatLocker returns a repository.SeatLocker backed by Redis.
func NewSeatLocker(locker *redislock.Locker) repository.SeatLocker {
	return &seatLocker{locker: locker}
}

func (s *seatLocker) AcquireSeatLock(ctx context.Context, cinemaID, seatNo, token string, ttl time.Duration) error {
	return s.locker.Acquire(ctx, redislock.SeatKey(cinemaID, seatNo), token, ttl)
}

func (s *seatLocker) ReleaseSeatLock(ctx context.Context, cinemaID, seatNo, token string) error {
	return s.locker.Release(ctx, redislock.SeatKey(cinemaID, seatNo), token)
}
