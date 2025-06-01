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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AttendanceRepository interface {
	CreateAttendance(ctx context.Context, attendance *models.Attendance) error
	GetAttendanceByID(ctx context.Context, id int64) (*models.Attendance, error)
	UpdateAttendance(ctx context.Context, attendance *models.Attendance) error
	DeleteAttendance(ctx context.Context, id int64) error
	ListAttendance(ctx context.Context, limit, offset int) ([]*models.Attendance, error)
	ListAttendanceWithFilters(ctx context.Context, studentID, disciplineID *int64, date *time.Time, limit, offset int) ([]*models.Attendance, error)
}

type AttendanceHandler struct {
	repo AttendanceRepository
}

func NewAttendanceHandler(repo AttendanceRepository) *AttendanceHandler {
	return &AttendanceHandler{repo: repo}
}

// @Summary Добавить посещаемость
// @Tags attendances
// @Accept json
// @Produce json
// @Param input body models.Attendance true "Посещаемость"
// @Success 201 {object} models.Attendance
// @Router /api/v1/attendances [post]
// @Security BearerAuth
func (h *AttendanceHandler) CreateAttendance(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.attendance_handler.CreateAttendance"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var a models.Attendance
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateAttendance(r.Context(), &a); err != nil {
			log.Error("failed to create attendance", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create attendance"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, a)
	}
}

// @Summary Получить посещаемость по ID
// @Tags attendances
// @Accept json
// @Produce json
// @Param id path int true "ID посещаемости"
// @Success 200 {object} models.Attendance
// @Router /api/v1/attendances/{id} [get]
// @Security BearerAuth
func (h *AttendanceHandler) GetAttendanceByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.attendance_handler.GetAttendanceByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid attendance id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid attendance id"))
			return
		}
		a, err := h.repo.GetAttendanceByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("attendance not found", slog.Int64("attendance_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("attendance not found"))
				return
			}
			log.Error("failed to get attendance", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get attendance"))
			return
		}
		render.JSON(w, r, a)
	}
}

// @Summary Обновить посещаемость
// @Tags attendances
// @Accept json
// @Produce json
// @Param id path int true "ID посещаемости"
// @Param input body models.Attendance true "Посещаемость"
// @Success 200 {object} models.Attendance
// @Router /api/v1/attendances/{id} [put]
// @Security BearerAuth
func (h *AttendanceHandler) UpdateAttendance(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.attendance_handler.UpdateAttendance"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid attendance id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid attendance id"))
			return
		}
		var a models.Attendance
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		a.AttendanceID = id
		if err := h.repo.UpdateAttendance(r.Context(), &a); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("attendance not found for update", slog.Int64("attendance_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("attendance not found"))
				return
			}
			log.Error("failed to update attendance", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update attendance"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, a)
	}
}

// @Summary Удалить посещаемость
// @Tags attendances
// @Accept json
// @Produce json
// @Param id path int true "ID посещаемости"
// @Success 204 {string} string "No Content"
// @Router /api/v1/attendances/{id} [delete]
// @Security BearerAuth
func (h *AttendanceHandler) DeleteAttendance(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.attendance_handler.DeleteAttendance"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid attendance id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid attendance id"))
			return
		}
		if err := h.repo.DeleteAttendance(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("attendance not found for delete", slog.Int64("attendance_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("attendance not found"))
				return
			}
			log.Error("failed to delete attendance", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete attendance"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список посещаемости с фильтрацией
// @Tags attendances
// @Accept json
// @Produce json
// @Param student_id query int false "ID студента"
// @Param discipline_id query int false "ID дисциплины"
// @Param from_date query string false "С даты (YYYY-MM-DD)"
// @Param to_date query string false "По дату (YYYY-MM-DD)"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Attendance
// @Router /api/v1/attendances [get]
// @Security BearerAuth
func (h *AttendanceHandler) ListAttendance(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.attendance_handler.ListAttendanceWithFilters"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var (
			studentID, disciplineID *int64
			date                    *time.Time
		)

		studentIDStr := r.URL.Query().Get("student_id")
		if studentIDStr != "" {
			id, err := strconv.ParseInt(studentIDStr, 10, 64)
			if err == nil {
				studentID = &id
			}
		}
		disciplineIDStr := r.URL.Query().Get("discipline_id")
		if disciplineIDStr != "" {
			id, err := strconv.ParseInt(disciplineIDStr, 10, 64)
			if err == nil {
				disciplineID = &id
			}
		}
		dateStr := r.URL.Query().Get("date")
		if dateStr != "" {
			parsed, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				date = &parsed
			}
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}

		items, err := h.repo.ListAttendanceWithFilters(r.Context(), studentID, disciplineID, date, limit, offset)
		if err != nil {
			log.Error("failed to list attendance with filters", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list attendance"))
			return
		}
		render.JSON(w, r, items)
	}
}
