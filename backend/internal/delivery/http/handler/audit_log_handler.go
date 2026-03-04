package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

type AuditLogHandler struct {
	auditUC domainusecase.AuditLogUseCase
}

func NewAuditLogHandler(auditUC domainusecase.AuditLogUseCase) *AuditLogHandler {
	return &AuditLogHandler{auditUC: auditUC}
}

// GET /api/admin/audit-logs
func (h *AuditLogHandler) List(c *gin.Context) {
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 20)

	result, err := h.auditUC.ListLogs(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
