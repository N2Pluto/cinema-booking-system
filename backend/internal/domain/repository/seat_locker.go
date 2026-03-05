package repository

import (
	"context"
	"time"
)

// SeatLocker is the domain-layer contract for a distributed seat lock.
// The concrete implementation lives in internal/repository/redis/.
type SeatLocker interface {
	// AcquireSeatLock tries to acquire an exclusive lock for seatNo in cinemaID.
	// token is a caller-owned identifier (e.g. userID) used to prevent
	// another caller from releasing a lock it does not own.
	// Returns ErrLockNotAcquired if the seat is already locked.
	AcquireSeatLock(ctx context.Context, cinemaID, seatNo, token string, ttl time.Duration) error

	// ReleaseSeatLock releases the lock for seatNo in cinemaID.
	// Only succeeds when the stored token matches, so a different user
	// cannot release someone else's lock.
	ReleaseSeatLock(ctx context.Context, cinemaID, seatNo, token string) error
}
