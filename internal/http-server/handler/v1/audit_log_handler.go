package v1

import (
	"context"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	resp "service/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AuditLogRepository interface {
	AddAuditLog(ctx context.Context, entry *models.AuditLog) error
	ListAuditLogs(ctx context.Context, limit, offset int) ([]*models.AuditLog, error)
}

type AuditLogHandler struct {
	repo AuditLogRepository
}

func NewAuditLogHandler(repo AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{repo: repo}
}

// @Summary Получить список аудитов
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.AuditLog
// @Router /api/v1/audit-logs [get]
// @Security BearerAuth
func (h *AuditLogHandler) ListAuditLogs(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.auditlog.ListAuditLogs"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}
		audits, err := h.repo.ListAuditLogs(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list audit logs", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list audit logs"))
			return
		}
		render.JSON(w, r, audits)
	}
}
