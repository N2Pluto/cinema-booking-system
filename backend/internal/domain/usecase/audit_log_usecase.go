package usecase

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type ListLogsResult struct {
	Data       []*entity.AuditLog `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

type AuditLogUseCase interface {
	ListLogs(ctx context.Context, page, limit int) (*ListLogsResult, error)
}
