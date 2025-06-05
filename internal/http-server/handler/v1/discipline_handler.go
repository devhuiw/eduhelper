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

type DisciplineRepository interface {
	CreateDiscipline(ctx context.Context, discipline *models.Discipline) error
	GetDisciplineByID(ctx context.Context, id int64) (*models.Discipline, error)
	UpdateDiscipline(ctx context.Context, discipline *models.Discipline) error
	DeleteDiscipline(ctx context.Context, id int64) error
	ListDiscipline(ctx context.Context, limit, offset int) ([]*models.Discipline, error)
	GetDisciplinePublicByID(ctx context.Context, id int64) (*models.DisciplinePublic, error)
	ListDisciplinePublic(ctx context.Context, limit, offset int, teacherID, studentGroupID, academicYearID *int64) ([]*models.DisciplinePublic, error)
}

type DisciplineHandler struct {
	repo      DisciplineRepository
	auditRepo AuditLogRepository
}

func NewDisciplineHandler(repo DisciplineRepository, auditRepo AuditLogRepository) *DisciplineHandler {
	return &DisciplineHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать дисциплину
// @Tags disciplines
// @Accept json
// @Produce json
// @Param input body models.Discipline true "Дисциплина"
// @Success 201 {object} models.Discipline
// @Router /api/v1/disciplines [post]
// @Security BearerAuth
func (h *DisciplineHandler) CreateDiscipline(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.CreateDiscipline"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var discipline models.Discipline
		if err := json.NewDecoder(r.Body).Decode(&discipline); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		if err := h.repo.CreateDiscipline(r.Context(), &discipline); err != nil {
			log.Error("failed to create discipline", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create discipline"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "discipline",
			RowID:      discipline.DisciplineID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(discipline),
			Comment:    utils.PtrToStr("Discipline created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, discipline)
	}
}

// @Summary Получить дисциплину по ID
// @Tags disciplines
// @Accept json
// @Produce json
// @Param id path int true "ID дисциплины"
// @Success 200 {object} models.Discipline
// @Router /api/v1/disciplines/{id} [get]
// @Security BearerAuth
func (h *DisciplineHandler) GetDisciplineByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.GetDisciplineByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid discipline id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid discipline id"))
			return
		}
		discipline, err := h.repo.GetDisciplineByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("discipline not found", slog.Int64("discipline_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("discipline not found"))
				return
			}
			log.Error("failed to get discipline", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get discipline"))
			return
		}
		render.JSON(w, r, discipline)
	}
}

// @Summary Обновить дисциплину
// @Tags disciplines
// @Accept json
// @Produce json
// @Param id path int true "ID дисциплины"
// @Param input body models.Discipline true "Дисциплина"
// @Success 200 {object} models.Discipline
// @Router /api/v1/disciplines/{id} [put]
// @Security BearerAuth
func (h *DisciplineHandler) UpdateDiscipline(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.UpdateDiscipline"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid discipline id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid discipline id"))
			return
		}
		var discipline models.Discipline
		if err := json.NewDecoder(r.Body).Decode(&discipline); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		discipline.DisciplineID = id
		oldData, _ := h.repo.GetDisciplineByID(r.Context(), id)
		if err := h.repo.UpdateDiscipline(r.Context(), &discipline); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("discipline not found for update", slog.Int64("discipline_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("discipline not found"))
				return
			}
			log.Error("failed to update discipline", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update discipline"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "discipline",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(discipline),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Discipline updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, discipline)
	}
}

// @Summary Удалить дисциплину
// @Tags disciplines
// @Accept json
// @Produce json
// @Param id path int true "ID дисциплины"
// @Success 204 {string} string "No Content"
// @Router /api/v1/disciplines/{id} [delete]
// @Security BearerAuth
func (h *DisciplineHandler) DeleteDiscipline(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.DeleteDiscipline"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid discipline id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid discipline id"))
			return
		}
		oldData, _ := h.repo.GetDisciplineByID(r.Context(), id)
		if err := h.repo.DeleteDiscipline(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("discipline not found for delete", slog.Int64("discipline_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("discipline not found"))
				return
			}
			log.Error("failed to delete discipline", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete discipline"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "discipline",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Discipline deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список дисциплин с фильтрацией
// @Tags disciplines
// @Accept json
// @Produce json
// @Param teacher_id query int false "ID преподавателя"
// @Param student_group_id query int false "ID группы"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Discipline
// @Router /api/v1/disciplines [get]
// @Security BearerAuth
func (h *DisciplineHandler) ListDiscipline(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.ListDiscipline"
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
		disciplines, err := h.repo.ListDiscipline(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list disciplines", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list disciplines"))
			return
		}
		render.JSON(w, r, disciplines)
	}
}

// @Summary Получить публичную дисциплину по ID
// @Tags disciplines
// @Accept json
// @Produce json
// @Param id path int true "ID дисциплины"
// @Success 200 {object} models.DisciplinePublic
// @Router /api/v1/disciplines/public/{id} [get]
// @Security BearerAuth
func (h *DisciplineHandler) GetDisciplinePublicByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.GetDisciplinePublicByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid discipline id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid discipline id"))
			return
		}
		discipline, err := h.repo.GetDisciplinePublicByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("discipline not found", slog.Int64("discipline_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("discipline not found"))
				return
			}
			log.Error("failed to get discipline public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get discipline public"))
			return
		}
		render.JSON(w, r, discipline)
	}
}

// @Summary Получить список публичных дисциплин
// @Tags disciplines
// @Accept json
// @Produce json
// @Param teacher_id query int false "ID преподавателя"
// @Param student_group_id query int false "ID группы студентов"
// @Param academic_year_id query int false "ID учебного года"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.DisciplinePublic
// @Router /api/v1/disciplines/public [get]
// @Security BearerAuth
func (h *DisciplineHandler) ListDisciplinePublic(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.discipline_handler.ListDisciplinePublic"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))
		if limit == 0 {
			limit = 20
		}

		// фильтры, можно расширить
		var (
			teacherID      *int64
			studentGroupID *int64
			academicYearID *int64
		)
		if val := q.Get("teacher_id"); val != "" {
			id, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				teacherID = &id
			}
		}
		if val := q.Get("student_group_id"); val != "" {
			id, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				studentGroupID = &id
			}
		}
		if val := q.Get("academic_year_id"); val != "" {
			id, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				academicYearID = &id
			}
		}

		disciplines, err := h.repo.ListDisciplinePublic(
			r.Context(), limit, offset, teacherID, studentGroupID, academicYearID,
		)
		if err != nil {
			log.Error("failed to list disciplines public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list disciplines public"))
			return
		}
		render.JSON(w, r, disciplines)
	}
}
