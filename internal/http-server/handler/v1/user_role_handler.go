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

func (h *UserRoleHandler) GetRolesByUserID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.userrole.GetRolesByUserID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid user id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid user id"))
			return
		}
		users_role, err := h.repo.GetRolesByUserID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user roles not found", slog.Any("users_role", users_role))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("user roles not found"))
				return
			}
			log.Error("failed to get user roles", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get user roles"))
			return
		}

		render.JSON(w, r, users_role)

	}
}
