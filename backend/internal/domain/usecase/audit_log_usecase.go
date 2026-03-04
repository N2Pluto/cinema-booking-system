package usecase

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type ListLogsResult struct {
	Data       []*entity.AuditLog
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type AuditLogUseCase interface {
	ListLogs(ctx context.Context, page, limit int) (*ListLogsResult, error)
}
