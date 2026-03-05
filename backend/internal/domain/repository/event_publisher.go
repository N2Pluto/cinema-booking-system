package repository

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

// EventPublisher is the domain-layer contract for publishing audit log events
// asynchronously via a message queue (Redis Pub-Sub).
// The concrete implementation lives in internal/repository/redis/.
type EventPublisher interface {
	PublishAuditLog(ctx context.Context, log *entity.AuditLog) error
}
