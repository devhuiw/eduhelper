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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SemesterRepository interface {
	CreateSemester(ctx context.Context, s *models.Semester) error
	GetSemesterByID(ctx context.Context, id int64) (*models.Semester, error)
	UpdateSemester(ctx context.Context, s *models.Semester) error
	DeleteSemester(ctx context.Context, id int64) error
	ListSemester(ctx context.Context, academicYearID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.Semester, error)
}

type SemesterHandler struct {
	repo      SemesterRepository
	auditRepo AuditLogRepository
}

func NewSemesterHandler(repo SemesterRepository, auditRepo AuditLogRepository) *SemesterHandler {
	return &SemesterHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать семестр
// @Tags semesters
// @Accept json
// @Produce json
// @Param input body models.Semester true "Семестр"
// @Success 201 {object} models.Semester
// @Router /api/v1/semesters [post]
// @Security BearerAuth
func (h *SemesterHandler) CreateSemester(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.semester_handler.CreateSemester"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var s models.Semester
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateSemester(r.Context(), &s); err != nil {
			log.Error("failed to create semester", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create semester"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "semestr",
			RowID:      s.SemesterID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(s),
			Comment:    utils.PtrToStr("Semestr created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, s)
	}
}

// @Summary Получить семестр по ID
// @Tags semesters
// @Accept json
// @Produce json
// @Param id path int true "ID семестра"
// @Success 200 {object} models.Semester
// @Router /api/v1/semesters/{id} [get]
// @Security BearerAuth
func (h *SemesterHandler) GetSemesterByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.semester_handler.GetSemesterByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid semester id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid semester id"))
			return
		}
		semester, err := h.repo.GetSemesterByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("semester not found", slog.Int64("semester_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("semester not found"))
				return
			}
			log.Error("failed to get semester", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get semester"))
			return
		}
		render.JSON(w, r, semester)
	}
}

// @Summary Обновить семестр
// @Tags semesters
// @Accept json
// @Produce json
// @Param id path int true "ID семестра"
// @Param input body models.Semester true "Семестр"
// @Success 200 {object} models.Semester
// @Router /api/v1/semesters/{id} [put]
// @Security BearerAuth
func (h *SemesterHandler) UpdateSemester(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.semester_handler.UpdateSemester"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid semester id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid semester id"))
			return
		}
		var s models.Semester
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		s.SemesterID = id
		oldData, _ := h.repo.GetSemesterByID(r.Context(), id)
		if err := h.repo.UpdateSemester(r.Context(), &s); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("semester not found for update", slog.Int64("semester_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("semester not found"))
				return
			}
			log.Error("failed to update semester", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update semester"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "semestr",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(s),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Semestr updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, s)
	}
}

// @Summary Удалить семестр
// @Tags semesters
// @Accept json
// @Produce json
// @Param id path int true "ID семестра"
// @Success 204 {string} string "No Content"
// @Router /api/v1/semesters/{id} [delete]
// @Security BearerAuth
func (h *SemesterHandler) DeleteSemester(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.semester_handler.DeleteSemester"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid semester id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid semester id"))
			return
		}
		oldData, _ := h.repo.GetSemesterByID(r.Context(), id)
		if err := h.repo.DeleteSemester(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("semester not found for delete", slog.Int64("semester_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("semester not found"))
				return
			}
			log.Error("failed to delete semester", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete semester"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "semestr",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Semestr deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список семестров с фильтрацией
// @Tags semesters
// @Accept json
// @Produce json
// @Param academic_year_id query int false "ID учебного года"
// @Param from_date query string false "С даты (YYYY-MM-DD)"
// @Param to_date query string false "По дату (YYYY-MM-DD)"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Semester
// @Router /api/v1/semesters [get]
// @Security BearerAuth
func (h *SemesterHandler) ListSemester(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.semester_handler.ListSemester"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var academicYearID *int64
		var fromDate, toDate *time.Time

		academicYearIDStr := r.URL.Query().Get("academic_year_id")
		if academicYearIDStr != "" {
			id, err := strconv.ParseInt(academicYearIDStr, 10, 64)
			if err == nil {
				academicYearID = &id
			}
		}
		fromDateStr := r.URL.Query().Get("from_date")
		if fromDateStr != "" {
			t, err := time.Parse("2006-01-02", fromDateStr)
			if err == nil {
				fromDate = &t
			}
		}
		toDateStr := r.URL.Query().Get("to_date")
		if toDateStr != "" {
			t, err := time.Parse("2006-01-02", toDateStr)
			if err == nil {
				toDate = &t
			}
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}

		semesters, err := h.repo.ListSemester(r.Context(), academicYearID, fromDate, toDate, limit, offset)
		if err != nil {
			log.Error("failed to list semesters", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list semesters"))
			return
		}
		render.JSON(w, r, semesters)
	}
}
