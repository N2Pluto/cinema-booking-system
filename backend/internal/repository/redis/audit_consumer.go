package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	goredis "github.com/redis/go-redis/v9"
)

// AuditConsumer subscribes to the Redis Pub-Sub channel "audit:events"
// and writes each received AuditLog to MongoDB asynchronously.
// This is the Message Queue consumer — it decouples audit writes from
// the hot booking path.
type AuditConsumer struct {
	rdb       *goredis.Client
	auditRepo repository.AuditLogRepository
}

// NewAuditConsumer returns an AuditConsumer ready to be started.
func NewAuditConsumer(rdb *goredis.Client, auditRepo repository.AuditLogRepository) *AuditConsumer {
	return &AuditConsumer{rdb: rdb, auditRepo: auditRepo}
}

// Start subscribes to the audit channel and processes messages in a goroutine
// until ctx is cancelled.
func (c *AuditConsumer) Start(ctx context.Context) {
	go func() {
		sub := c.rdb.Subscribe(ctx, auditLogChannel)
		defer sub.Close()

		ch := sub.Channel()
		log.Println("audit_consumer: listening on channel", auditLogChannel)

		for {
			select {
			case <-ctx.Done():
				log.Println("audit_consumer: shutting down")
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				c.handle(ctx, msg.Payload)
			}
		}
	}()
}

func (c *AuditConsumer) handle(ctx context.Context, payload string) {
	var auditLog entity.AuditLog
	if err := json.Unmarshal([]byte(payload), &auditLog); err != nil {
		log.Printf("audit_consumer: unmarshal: %v", err)
		return
	}
	if err := c.auditRepo.Create(ctx, &auditLog); err != nil {
		log.Printf("audit_consumer: create audit log: %v", err)
	}
}
