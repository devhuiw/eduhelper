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

type StudentRepository interface {
	CreateStudent(ctx context.Context, student *models.Student) error
	GetStudentByID(ctx context.Context, userID int64) (*models.Student, error)
	GetStudentPublicByID(ctx context.Context, userID int64) (*models.StudentPublic, error)
	UpdateStudent(ctx context.Context, student *models.Student) error
	DeleteStudent(ctx context.Context, userID int64) error
	ListStudent(ctx context.Context, limit, offset int) ([]*models.Student, error)
	ListStudentPublic(ctx context.Context, limit, offset int) ([]*models.StudentPublic, error)
}

type StudentHandler struct {
	repo StudentRepository
}

func NewStudentHandler(repo StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func (h *StudentHandler) CreateStudent(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.CreateStudent"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var student models.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateStudent(r.Context(), &student); err != nil {
			log.Error("failed to create student", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create student"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, student)
	}
}

func (h *StudentHandler) GetStudentByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.GetStudentByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid student id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid student id"))
			return
		}
		student, err := h.repo.GetStudentByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student not found", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student not found"))
				return
			}
			log.Error("failed to get student", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get student"))
			return
		}
		render.JSON(w, r, student)
	}
}

func (h *StudentHandler) GetStudentPublicByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.GetStudentPublicByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid student id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid student id"))
			return
		}
		student, err := h.repo.GetStudentPublicByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student not found", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student not found"))
				return
			}
			log.Error("failed to get student public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get student public"))
			return
		}
		render.JSON(w, r, student)
	}
}

func (h *StudentHandler) UpdateStudent(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.UpdateStudent"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid student id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid student id"))
			return
		}
		var student models.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		student.UserID = id
		if err := h.repo.UpdateStudent(r.Context(), &student); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student not found for update", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student not found"))
				return
			}
			log.Error("failed to update student", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update student"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, student)
	}
}

func (h *StudentHandler) DeleteStudent(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.DeleteStudent"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid student id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid student id"))
			return
		}
		if err := h.repo.DeleteStudent(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student not found for delete", slog.Int64("user_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student not found"))
				return
			}
			log.Error("failed to delete student", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete student"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *StudentHandler) ListStudent(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.ListStudent"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}
		students, err := h.repo.ListStudent(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list students", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list students"))
			return
		}
		render.JSON(w, r, students)
	}
}

func (h *StudentHandler) ListStudentPublic(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.student_handler.ListStudentPublic"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}
		students, err := h.repo.ListStudentPublic(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list students public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list students public"))
			return
		}
		render.JSON(w, r, students)
	}
}
