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
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type RolePermissionRepository interface {
	AssignPermission(ctx context.Context, roleID, permissionID int64) error
	RemovePermission(ctx context.Context, roleID, permissionID int64) error
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]*models.Permission, error)
}

type RolePermissionHandler struct {
	repo RolePermissionRepository
}

func NewRolePermissionHandler(repo RolePermissionRepository) *RolePermissionHandler {
	return &RolePermissionHandler{repo: repo}
}

type assignPermissionInput struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

// @Summary Назначить право роли
// @Tags role-permissions
// @Accept json
// @Produce json
// @Param input body assignPermissionInput true "Роль и право"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/role-permissions/assign [post]
// @Security BearerAuth
func (h *RolePermissionHandler) AssignPermission(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.rolepermission.AssignPermission"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var input assignPermissionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.AssignPermission(r.Context(), input.RoleID, input.PermissionID); err != nil {
			log.Error("failed to assign permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to assign permission"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}

// @Summary Удалить право у роли
// @Tags role-permissions
// @Accept json
// @Produce json
// @Param input body assignPermissionInput true "Роль и право"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/role-permissions/remove [post]
// @Security BearerAuth
func (h *RolePermissionHandler) RemovePermission(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.rolepermission.RemovePermission"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var input assignPermissionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.RemovePermission(r.Context(), input.RoleID, input.PermissionID); err != nil {
			log.Error("failed to remove permission", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to remove permission"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}

// @Summary Получить список прав роли
// @Tags role-permissions
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {array} models.Permission
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/role-permissions/{id} [get]
// @Security BearerAuth
func (h RolePermissionHandler) GetPermissionsByRoleID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.rolepermission.GetPermissionsByRoleID"
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		role_id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid role id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid role id"))
			return
		}
		permissions, err := h.repo.GetPermissionsByRoleID(r.Context(), role_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("permissions for role id not found", slog.Any("permissions", permissions))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("permissions for role id not found"))
				return
			}
			log.Error("failed to get permissions for role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get permissions for role"))
			return
		}

		render.JSON(w, r, permissions)
	}
}
