package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	goredis "github.com/redis/go-redis/v9"
)

const auditLogChannel = "audit:events"

type eventPublisher struct {
	rdb *goredis.Client
}

// NewEventPublisher returns a repository.EventPublisher that publishes
// audit log events to the Redis Pub-Sub channel "audit:events".
func NewEventPublisher(rdb *goredis.Client) repository.EventPublisher {
	return &eventPublisher{rdb: rdb}
}

func (p *eventPublisher) PublishAuditLog(ctx context.Context, log *entity.AuditLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("event_publisher: marshal: %w", err)
	}
	return p.rdb.Publish(ctx, auditLogChannel, data).Err()
}
