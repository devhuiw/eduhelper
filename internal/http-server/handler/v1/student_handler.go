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
	repo      StudentRepository
	auditRepo AuditLogRepository
}

func NewStudentHandler(repo StudentRepository, auditRepo AuditLogRepository) *StudentHandler {
	return &StudentHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать студента
// @Tags students
// @Accept json
// @Produce json
// @Param input body models.Student true "Студент"
// @Success 201 {object} models.Student
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/students [post]
// @Security BearerAuth
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student",
			RowID:      student.UserID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(student),
			Comment:    utils.PtrToStr("Student created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, student)
	}
}

// @Summary Получить студента по ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "ID студента"
// @Success 200 {object} models.Student
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/students/{id} [get]
// @Security BearerAuth
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

// @Summary Получить публичного студента по ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "ID студента"
// @Success 200 {object} models.StudentPublic
// @Failure 404 {object} resp.Response
// @Router /api/v1/students/public/{id} [get]
// @Security BearerAuth
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

// @Summary Обновить данные студента
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "ID студента"
// @Param input body models.Student true "Студент"
// @Success 200 {object} models.Student
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/students/{id} [put]
// @Security BearerAuth
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
		oldData, _ := h.repo.GetStudentByID(r.Context(), id)
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(student),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Student updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, student)
	}
}

// @Summary Удалить студента
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "ID студента"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/students/{id} [delete]
// @Security BearerAuth
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
		oldData, _ := h.repo.GetStudentByID(r.Context(), id)
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Student deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список студентов
// @Tags students
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Student
// @Failure 500 {object} resp.Response
// @Router /api/v1/students [get]
// @Security BearerAuth
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

// @Summary Получить публичный список студентов
// @Tags students
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.StudentPublic
// @Router /api/v1/students/public [get]
// @Security BearerAuth
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
