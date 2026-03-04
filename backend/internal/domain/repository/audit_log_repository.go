package repository

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type AuditLogFilter struct {
	Page  int
	Limit int
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	List(ctx context.Context, f AuditLogFilter) ([]*entity.AuditLog, int64, error)
}
