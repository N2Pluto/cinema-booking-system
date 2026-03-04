package audit_log

import (
	"context"
	"fmt"
	"math"

	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

var _ domainusecase.AuditLogUseCase = (*UseCase)(nil)

type UseCase struct {
	auditRepo repository.AuditLogRepository
}

func NewUseCase(auditRepo repository.AuditLogRepository) *UseCase {
	return &UseCase{auditRepo: auditRepo}
}

func (uc *UseCase) ListLogs(ctx context.Context, page, limit int) (*domainusecase.ListLogsResult, error) {
	if limit < 1 {
		limit = 20
	}

	logs, total, err := uc.auditRepo.List(ctx, repository.AuditLogFilter{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("list logs: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domainusecase.ListLogsResult{
		Data:       logs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
