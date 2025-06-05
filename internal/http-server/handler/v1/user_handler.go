package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"service/internal/domain/models"

	resp "service/internal/lib/api/response"
	"service/internal/lib/utils"
	"strconv"

	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserRepository interface {
	CreateClient(ctx context.Context, user *models.User) error
	GetClientByID(ctx context.Context, id int64) (*models.User, error)
	GetClientByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateClient(ctx context.Context, user *models.User) error
	DeleteClient(ctx context.Context, id int64) error
	ListClient(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type UserHandler struct {
	repo      UserRepository
	auditRepo AuditLogRepository
}

func NewUserHandler(repo UserRepository, auditRepo AuditLogRepository) *UserHandler {
	return &UserHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.User true "Пользователь"
// @Success 201 {object} models.User
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.user_handler.CreateUser"
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
		if err := h.repo.CreateClient(r.Context(), &user); err != nil {
			log.Error("failed to create user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create user"))
			return
		}

		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "user",
			RowID:      user.UserID,
			ActionType: "INSERT",
			NewData:    utils.PtrToJSON(user),
			Comment:    utils.PtrToStr("User created"),
		})

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, user)
	}
}

// @Summary Получить пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
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
		user, err := h.repo.GetClientByID(r.Context(), id)
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

// @Summary Обновить пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body models.User true "Пользователь"
// @Success 200 {object} models.User
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
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
		oldUser, _ := h.repo.GetClientByID(r.Context(), id)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		user.UserID = id
		if err := h.repo.UpdateClient(r.Context(), &user); err != nil {
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

		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "user",
			RowID:      user.UserID,
			ActionType: "UPDATE",
			OldData:    utils.PtrToJSON(oldUser),
			NewData:    utils.PtrToJSON(user),
			Comment:    utils.PtrToStr("User updated"),
		})

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, user)
	}
}

// @Summary Удалить пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
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
		oldUser, _ := h.repo.GetClientByID(r.Context(), id)
		if err := h.repo.DeleteClient(r.Context(), id); err != nil {
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

		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "user",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldUser),
			Comment:    utils.PtrToStr("User deleted"),
		})

		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список пользователей
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.User
// @Failure 500 {object} resp.Response
// @Router /api/v1/users [get]
// @Security BearerAuth
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
		users, err := h.repo.ListClient(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list users", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list users"))
			return
		}
		render.JSON(w, r, users)
	}
}
