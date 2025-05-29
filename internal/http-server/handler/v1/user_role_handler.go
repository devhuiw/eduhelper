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

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]*models.UserRole, error)
}

type UserRoleHandler struct {
	repo UserRoleRepository
}

func NewUserRoleHandler(repo UserRoleRepository) *UserRoleHandler {
	return &UserRoleHandler{repo: repo}
}

type assignRoleInput struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// POST /api/v1/user-roles/assign
func (h *UserRoleHandler) AssignRole(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.userrole.AssignRole"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var input assignRoleInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.AssignRole(r.Context(), input.UserID, input.RoleID); err != nil {
			log.Error("failed to assign role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to assign role"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}

// POST /api/v1/user-roles/remove
func (h *UserRoleHandler) RemoveRole(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.userrole.RemoveRole"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var input assignRoleInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.RemoveRole(r.Context(), input.UserID, input.RoleID); err != nil {
			log.Error("failed to remove role", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to remove role"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}
