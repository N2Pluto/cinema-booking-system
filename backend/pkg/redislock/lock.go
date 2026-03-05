// Package redislock provides a Redis-backed distributed lock.
//
// Strategy:
//   - Acquire: SET lock:{key} {token} NX EX {ttl}
//     NX ensures only one caller succeeds for the same key.
//     The token (caller-owned UUID) prevents a different caller from
//     releasing a lock it does not own.
//   - Release: Lua script — DEL the key only when its value equals our token.
//     This is atomic; no race between GET and DEL.
//
// Seat locking uses key  lock:seat:{cinemaID}:{seatNo}  with TTL = 5 minutes.
// If a user does not confirm within 5 minutes the key expires automatically,
// and the background ticker also explicitly cleans up the DB state.
package redislock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrNotAcquired is returned when the lock key is already held by someone else.
var ErrNotAcquired = errors.New("lock: not acquired — already held")

// releaseScript releases a lock only when the stored token matches the caller's token.
// Returns 1 on success, 0 if the key was already gone or owned by another caller.
var releaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`)

// Locker holds a reference to the Redis client.
type Locker struct {
	rdb *redis.Client
}

// New returns a new Locker backed by the given Redis client.
func New(rdb *redis.Client) *Locker {
	return &Locker{rdb: rdb}
}

// Acquire tries to set lock:{key} = token with NX EX ttl.
// Returns (token, nil) on success, ("", ErrNotAcquired) if already locked.
func (l *Locker) Acquire(ctx context.Context, key, token string, ttl time.Duration) error {
	ok, err := l.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAcquired
	}
	return nil
}

// Release deletes lock:{key} only if the stored value equals token.
func (l *Locker) Release(ctx context.Context, key, token string) error {
	return releaseScript.Run(ctx, l.rdb, []string{key}, token).Err()
}

// SeatKey returns the canonical Redis key for a cinema seat lock.
func SeatKey(cinemaID, seatNo string) string {
	return "lock:seat:" + cinemaID + ":" + seatNo
}
