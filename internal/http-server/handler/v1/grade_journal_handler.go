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

type GradeJournalRepository interface {
	CreateGradeJournal(ctx context.Context, g *models.GradeJournal) error
	GetGradeJournalByID(ctx context.Context, id int64) (*models.GradeJournal, error)
	UpdateGradeJournal(ctx context.Context, g *models.GradeJournal) error
	DeleteGradeJournal(ctx context.Context, id int64) error
	ListGradeJournal(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.GradeJournal, error)
	ListGradeJournalPublic(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.GradeJournalPublic, error)
	GetAverageGrade(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time) (float64, error)
}

type GradeJournalHandler struct {
	repo      GradeJournalRepository
	auditRepo AuditLogRepository
}

func NewGradeJournalHandler(repo GradeJournalRepository, auditRepo AuditLogRepository) *GradeJournalHandler {
	return &GradeJournalHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Добавить запись в журнал оценок
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param input body models.GradeJournal true "Запись"
// @Success 201 {object} models.GradeJournal
// @Router /api/v1/gradejournals [post]
// @Security BearerAuth
func (h *GradeJournalHandler) CreateGradeJournal(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.CreateGradeJournal"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var g models.GradeJournal
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateGradeJournal(r.Context(), &g); err != nil {
			log.Error("failed to create gradejournal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create gradejournal"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "grade_journal",
			RowID:      g.GradeJournalID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(g),
			Comment:    utils.PtrToStr("Grade_Journal created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, g)
	}
}

// @Summary Получить запись журнала по ID
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} models.GradeJournal
// @Router /api/v1/gradejournals/{id} [get]
// @Security BearerAuth
func (h *GradeJournalHandler) GetGradeJournalByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.GetGradeJournalByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid gradejournal id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid gradejournal id"))
			return
		}
		g, err := h.repo.GetGradeJournalByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("gradejournal not found", slog.Int64("gradejournal_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("gradejournal not found"))
				return
			}
			log.Error("failed to get gradejournal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get gradejournal"))
			return
		}
		render.JSON(w, r, g)
	}
}

// @Summary Обновить запись в журнале
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Param input body models.GradeJournal true "Запись"
// @Success 200 {object} models.GradeJournal
// @Router /api/v1/gradejournals/{id} [put]
// @Security BearerAuth
func (h *GradeJournalHandler) UpdateGradeJournal(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.UpdateGradeJournal"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid gradejournal id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid gradejournal id"))
			return
		}
		var g models.GradeJournal
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		g.GradeJournalID = id
		oldData, _ := h.repo.GetGradeJournalByID(r.Context(), id)
		if err := h.repo.UpdateGradeJournal(r.Context(), &g); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("gradejournal not found for update", slog.Int64("gradejournal_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("gradejournal not found"))
				return
			}
			log.Error("failed to update gradejournal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update gradejournal"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "grade_journal",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(g),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Grade_Journal created"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, g)
	}
}

// @Summary Удалить запись из журнала
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Success 204 {string} string "No Content"
// @Router /api/v1/gradejournals/{id} [delete]
// @Security BearerAuth
func (h *GradeJournalHandler) DeleteGradeJournal(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.DeleteGradeJournal"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid gradejournal id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid gradejournal id"))
			return
		}
		oldData, _ := h.repo.GetGradeJournalByID(r.Context(), id)
		if err := h.repo.DeleteGradeJournal(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("gradejournal not found for delete", slog.Int64("gradejournal_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("gradejournal not found"))
				return
			}
			log.Error("failed to delete gradejournal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete gradejournal"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "grade_journal",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Grade_Journal deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список оценок с фильтрацией
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param student_id query int false "ID студента"
// @Param discipline_id query int false "ID дисциплины"
// @Param from_date query string false "С даты (YYYY-MM-DD)"
// @Param to_date query string false "По дату (YYYY-MM-DD)"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.GradeJournal
// @Router /api/v1/gradejournals [get]
// @Security BearerAuth
func (h *GradeJournalHandler) ListGradeJournal(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.ListGradeJournal"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var studentID, disciplineID *int64
		var fromDate, toDate *time.Time

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
		fromDateStr := r.URL.Query().Get("from_date")
		if fromDateStr != "" {
			d, err := time.Parse("2006-01-02", fromDateStr)
			if err == nil {
				fromDate = &d
			}
		}
		toDateStr := r.URL.Query().Get("to_date")
		if toDateStr != "" {
			d, err := time.Parse("2006-01-02", toDateStr)
			if err == nil {
				toDate = &d
			}
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}

		items, err := h.repo.ListGradeJournal(r.Context(), studentID, disciplineID, fromDate, toDate, limit, offset)
		if err != nil {
			log.Error("failed to list gradejournals", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list gradejournals"))
			return
		}
		render.JSON(w, r, items)
	}
}

// @Summary Получить список публичных оценок
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param student_id query int false "ID студента"
// @Param discipline_id query int false "ID дисциплины"
// @Param from_date query string false "С даты (YYYY-MM-DD)"
// @Param to_date query string false "По дату (YYYY-MM-DD)"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.GradeJournalPublic
// @Router /api/v1/gradejournals/public [get]
// @Security BearerAuth
func (h *GradeJournalHandler) ListGradeJournalPublic(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.ListGradeJournalPublic"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var studentID, disciplineID *int64
		var fromDate, toDate *time.Time

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
		fromDateStr := r.URL.Query().Get("from_date")
		if fromDateStr != "" {
			d, err := time.Parse("2006-01-02", fromDateStr)
			if err == nil {
				fromDate = &d
			}
		}
		toDateStr := r.URL.Query().Get("to_date")
		if toDateStr != "" {
			d, err := time.Parse("2006-01-02", toDateStr)
			if err == nil {
				toDate = &d
			}
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}

		items, err := h.repo.ListGradeJournalPublic(r.Context(), studentID, disciplineID, fromDate, toDate, limit, offset)
		if err != nil {
			log.Error("failed to list gradejournals public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list gradejournals public"))
			return
		}
		render.JSON(w, r, items)
	}
}

// @Summary Получить средний балл
// @Tags gradejournals
// @Accept json
// @Produce json
// @Param student_id query int false "ID студента"
// @Param discipline_id query int false "ID дисциплины"
// @Param from_date query string false "С даты (YYYY-MM-DD)"
// @Param to_date query string false "По дату (YYYY-MM-DD)"
// @Success 200 {object} map[string]float64
// @Router /api/v1/gradejournals/average [get]
// @Security BearerAuth
func (h *GradeJournalHandler) GetAverageGrade(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.gradejournal_handler.GetAverageGrade"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var studentID, disciplineID *int64
		var fromDate, toDate *time.Time

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
		fromDateStr := r.URL.Query().Get("from_date")
		if fromDateStr != "" {
			d, err := time.Parse("2006-01-02", fromDateStr)
			if err == nil {
				fromDate = &d
			}
		}
		toDateStr := r.URL.Query().Get("to_date")
		if toDateStr != "" {
			d, err := time.Parse("2006-01-02", toDateStr)
			if err == nil {
				toDate = &d
			}
		}

		avg, err := h.repo.GetAverageGrade(r.Context(), studentID, disciplineID, fromDate, toDate)
		if err != nil {
			log.Error("failed to get average grade", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get average grade"))
			return
		}
		render.JSON(w, r, map[string]float64{"average_grade": avg})
	}
}
