package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	resp "service/internal/lib/api/response"

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

// POST /api/v1/role-permissions/assign
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

// POST /api/v1/role-permissions/remove
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
