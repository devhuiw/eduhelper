package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	resp "service/internal/lib/api/response"
	"strconv"

	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// POST /api/v1/users
func (h *UserHandler) CreateUser(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user.CreateUser"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.Create(r.Context(), &user); err != nil {
			log.Error("failed to create user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create user"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, user)
	}
}

// GET /api/v1/users/{id}
func (h *UserHandler) GetUserByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user.GetUserByID"
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
		user, err := h.repo.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user not found", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("user not found"))
				return
			}
			log.Error("failed to get user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get user"))
			return
		}
		render.JSON(w, r, user)
	}
}

// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user.UpdateUser"
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
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		user.UserID = id
		if err := h.repo.Update(r.Context(), &user); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user not found for update", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("user not found"))
				return
			}
			log.Info("failed to update user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update user"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, user)
	}
}

// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user.DeleteUser"
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
		if err := h.repo.Delete(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user not found for delete", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("user not found"))
				return
			}
			log.Error("failed to delete user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete user"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// GET /api/v1/users?limit=10&offset=0
func (h *UserHandler) ListUsers(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user.ListUsers"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)
		if limit == 0 {
			limit = 20
		}
		users, err := h.repo.List(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list users", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list users"))
			return
		}
		render.JSON(w, r, users)
	}
}
