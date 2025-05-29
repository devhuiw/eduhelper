package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	ware "service/internal/http-server/middleware"
	resp "service/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type TeacherRepository interface {
	CreateTeacher(ctx context.Context, teacher *models.Teacher) error
	GetTeacherByID(ctx context.Context, userID int64) (*models.Teacher, error)
	GetTeacherPublicByID(ctx context.Context, userID int64) (*models.TeacherPublic, error)
	UpdateTeacher(ctx context.Context, teacher *models.Teacher) error
	DeleteTeacher(ctx context.Context, userID int64) error
	ListTeacher(ctx context.Context, limit, offset int) ([]*models.Teacher, error)
	ListTeacherPublic(ctx context.Context, limit, offset int) ([]*models.TeacherPublic, error)
}

type TeacherHandler struct {
	repo TeacherRepository
}

func NewTeacherHandler(repo TeacherRepository) *TeacherHandler {
	return &TeacherHandler{repo: repo}
}

func (h *TeacherHandler) CreateTeacher(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.Create"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var teacher models.Teacher
		if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateTeacher(r.Context(), &teacher); err != nil {
			log.Error("failed to create teacher", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create teacher"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) GetTeacherByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.GetTeacherByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid teacher id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid teacher id"))
			return
		}
		teacher, err := h.repo.GetTeacherByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("teacher not found", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Error("failed to get teacher", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get teacher"))
			return
		}
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) GetTeacherPublicByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.GetTeacherPublicByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid teacher id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid teacher id"))
			return
		}
		teacher, err := h.repo.GetTeacherPublicByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("teacher not found", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Error("failed to get teacher", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get teacher"))
			return
		}
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) GetMyTeacherProfile(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.GetMyTeacherProfile"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		claims := ware.GetUserClaims(r)
		teacherId := claims["id"].(int64)
		teacher, err := h.repo.GetTeacherByID(r.Context(), teacherId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("teacher not found", slog.Int64("user_id", teacherId))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Error("failed to get teacher", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get teacher"))
			return
		}
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) UpdateTeacher(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.UpdateTeacher"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		teacherId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid teacher id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		var teacher models.Teacher
		if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		teacher.UserID = teacherId
		if err := h.repo.UpdateTeacher(r.Context(), &teacher); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user not found for update", slog.Int64("user_id", teacherId))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Info("failed to update user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) UpdateMyTeacherProfile(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.UpdateMyTeacherProfile"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		claims := ware.GetUserClaims(r)
		teacherId := claims["id"].(int64)
		var teacher models.Teacher
		if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		teacher.UserID = teacherId
		if err := h.repo.UpdateTeacher(r.Context(), &teacher); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("user not found for update", slog.Int64("user_id", teacherId))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Info("failed to update user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, teacher)
	}
}

func (h *TeacherHandler) DeleteTeacher(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher.DeleteTeacher"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid teacher id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid teacher id"))
			return
		}
		if err := h.repo.DeleteTeacher(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("teacher not found for delete", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("teacher not found"))
				return
			}
			log.Error("failed to delete teacher", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete teacher"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *TeacherHandler) ListTeacher(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.ListTeacher"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		teachers, err := h.repo.ListTeacher(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list teachers", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list teachers"))
			return
		}
		render.JSON(w, r, teachers)
	}
}

func (h *TeacherHandler) ListTeacherPublic(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.teacher_handler.ListTeacherPublic"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		teachers, err := h.repo.ListTeacherPublic(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list public teachers", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list public teachers"))
			return
		}
		render.JSON(w, r, teachers)
	}
}
