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

type RoleRepository interface {
	CreateRole(ctx context.Context, role *models.Role) (int64, error)
	GetRoleByID(ctx context.Context, id int64) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id int64) error
	ListRole(ctx context.Context) ([]*models.Role, error)
}

type RoleHandler struct {
	repo      RoleRepository
	auditRepo AuditLogRepository
}

func NewRoleHandler(repo RoleRepository, auditRepo AuditLogRepository) *RoleHandler {
	return &RoleHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать роль
// @Tags roles
// @Accept json
// @Produce json
// @Param input body models.Role true "Роль"
// @Success 201 {object} models.Role
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/roles [post]
// @Security BearerAuth
func (h *RoleHandler) CreateRole(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.role.CreateRole"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var role models.Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		id, err := h.repo.CreateRole(r.Context(), &role)
		if err != nil {
			log.Error("failed to create role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create role"))
			return
		}
		role.RoleID = id
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "Role",
			RowID:      role.RoleID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(role),
			Comment:    utils.PtrToStr("Role created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, role)
	}
}

// @Summary Получить роль по ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {object} models.Role
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/roles/{id} [get]
// @Security BearerAuth
func (h *RoleHandler) GetRoleByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.role.GetRoleByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid role id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid role id"))
			return
		}
		role, err := h.repo.GetRoleByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("role not found", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("role not found"))
				return
			}
			log.Error("failed to get role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get role"))
			return
		}
		render.JSON(w, r, role)
	}
}

// @Summary Обновить роль
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Param input body models.Role true "Роль"
// @Success 200 {object} models.Role
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/roles/{id} [put]
// @Security BearerAuth
func (h *RoleHandler) UpdateRole(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.role.UpdateRole"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid role id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		var role models.Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		role.RoleID = id
		oldData, _ := h.repo.GetRoleByID(r.Context(), id)
		if err := h.repo.UpdateRole(r.Context(), &role); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("role not found for update", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("role not found"))
				return
			}
			log.Error("failed to update role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update role"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "Role",
			RowID:      role.RoleID,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(role),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Role updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, role)
	}
}

// @Summary Удалить роль
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/roles/{id} [delete]
// @Security BearerAuth
func (h *RoleHandler) DeleteRole(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.role.DeleteRole"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid role id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid role id"))
			return
		}
		oldData, _ := h.repo.GetRoleByID(r.Context(), id)
		if err := h.repo.DeleteRole(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("role not found for delete", slog.Int64("id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("role not found"))
				return
			}
			log.Error("failed to delete role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete role"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "Role",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Role deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список ролей
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} models.Role
// @Failure 500 {object} resp.Response
// @Router /api/v1/roles [get]
// @Security BearerAuth
func (h *RoleHandler) ListRoles(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.role.ListRoles"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		roles, err := h.repo.ListRole(r.Context())
		if err != nil {
			log.Error("failed to list roles", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list roles"))
			return
		}
		render.JSON(w, r, roles)
	}
}
