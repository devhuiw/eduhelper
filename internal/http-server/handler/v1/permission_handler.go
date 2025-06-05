package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	resp "service/internal/lib/api/response"
	"service/internal/lib/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type PermissionRepository interface {
	CreatePermission(ctx context.Context, perm *models.Permission) error
	GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	UpdatePermission(ctx context.Context, perm *models.Permission) error
	DeletePermission(ctx context.Context, id int64) error
	ListPermission(ctx context.Context, limit, offset int) ([]*models.Permission, error)
}

type PermissionHandler struct {
	repo      PermissionRepository
	auditRepo AuditLogRepository
}

func NewPermissionHandler(repo PermissionRepository, auditRepo AuditLogRepository) *PermissionHandler {
	return &PermissionHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать право
// @Tags permissions
// @Accept json
// @Produce json
// @Param input body models.Permission true "Право"
// @Success 201 {object} models.Permission
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/permissions [post]
// @Security BearerAuth
func (h *PermissionHandler) CreatePermission(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.permission.CreatePermission"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var perm models.Permission
		if err := json.NewDecoder(r.Body).Decode(&perm); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreatePermission(r.Context(), &perm); err != nil {
			log.Error("failed to create permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create permission"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "permission",
			RowID:      perm.PermissionID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(perm),
			Comment:    utils.PtrToStr("Permission created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, perm)
	}
}

// @Summary Получить право по ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "ID права"
// @Success 200 {object} models.Permission
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/permissions/{id} [get]
// @Security BearerAuth
func (h *PermissionHandler) GetPermissionByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.permission.GetPermissionByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid permission id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid permission id"))
			return
		}
		perm, err := h.repo.GetPermissionByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("permission not found", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("permission not found"))
				return
			}
			log.Error("failed to get permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get permission"))
			return
		}
		render.JSON(w, r, perm)
	}
}

// @Summary Обновить право
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "ID права"
// @Param input body models.Permission true "Право"
// @Success 200 {object} models.Permission
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/permissions/{id} [put]
// @Security BearerAuth
func (h *PermissionHandler) UpdatePermission(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.permission.UpdatePermission"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid permission id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		var perm models.Permission
		if err := json.NewDecoder(r.Body).Decode(&perm); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		perm.PermissionID = id
		oldData, _ := h.repo.GetPermissionByID(r.Context(), id)
		if err := h.repo.UpdatePermission(r.Context(), &perm); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("permission not found for update", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("permission not found"))
				return
			}
			log.Error("failed to update permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update permission"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "permission",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(perm),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Permission updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, perm)
	}
}

// @Summary Удалить право
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "ID права"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/permissions/{id} [delete]
// @Security BearerAuth
func (h *PermissionHandler) DeletePermission(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.permission.DeletePermission"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid permission id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid permission id"))
			return
		}
		oldData, _ := h.repo.GetPermissionByID(r.Context(), id)
		if err := h.repo.DeletePermission(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("permission not found for delete", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("permission not found"))
				return
			}
			log.Error("failed to delete permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete permission"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "permission",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Permission created"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список прав
// @Tags permissions
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Permission
// @Failure 500 {object} resp.Response
// @Router /api/v1/permissions [get]
// @Security BearerAuth
func (h *PermissionHandler) ListPermissions(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.permission.ListPermissions"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}
		perms, err := h.repo.ListPermission(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list permissions", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list permissions"))
			return
		}
		render.JSON(w, r, perms)
	}
}
