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

type CurriculumRepository interface {
	CreateCurriculum(ctx context.Context, c *models.Curriculum) error
	GetCurriculumByID(ctx context.Context, id int64) (*models.Curriculum, error)
	UpdateCurriculum(ctx context.Context, c *models.Curriculum) error
	DeleteCurriculum(ctx context.Context, id int64) error
	ListCurriculum(ctx context.Context, semesterID, disciplineID *int64, limit, offset int) ([]*models.Curriculum, error)
}

type CurriculumHandler struct {
	repo CurriculumRepository
}

func NewCurriculumHandler(repo CurriculumRepository) *CurriculumHandler {
	return &CurriculumHandler{repo: repo}
}

// @Summary Создать учебный план
// @Tags curriculums
// @Accept json
// @Produce json
// @Param input body models.Curriculum true "Учебный план"
// @Success 201 {object} models.Curriculum
// @Router /api/v1/curriculums [post]
// @Security BearerAuth
func (h *CurriculumHandler) CreateCurriculum(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.curriculum_handler.CreateCurriculum"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var c models.Curriculum
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateCurriculum(r.Context(), &c); err != nil {
			log.Error("failed to create curriculum", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create curriculum"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, c)
	}
}

// @Summary Получить учебный план по ID
// @Tags curriculums
// @Accept json
// @Produce json
// @Param id path int true "ID учебного плана"
// @Success 200 {object} models.Curriculum
// @Router /api/v1/curriculums/{id} [get]
// @Security BearerAuth
func (h *CurriculumHandler) GetCurriculumByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.curriculum_handler.GetCurriculumByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid curriculum id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid curriculum id"))
			return
		}
		c, err := h.repo.GetCurriculumByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("curriculum not found", slog.Int64("curriculum_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("curriculum not found"))
				return
			}
			log.Error("failed to get curriculum", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get curriculum"))
			return
		}
		render.JSON(w, r, c)
	}
}

// @Summary Обновить учебный план
// @Tags curriculums
// @Accept json
// @Produce json
// @Param id path int true "ID учебного плана"
// @Param input body models.Curriculum true "Учебный план"
// @Success 200 {object} models.Curriculum
// @Router /api/v1/curriculums/{id} [put]
// @Security BearerAuth
func (h *CurriculumHandler) UpdateCurriculum(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.curriculum_handler.UpdateCurriculum"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid curriculum id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid curriculum id"))
			return
		}
		var c models.Curriculum
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		c.CurriculumID = id
		if err := h.repo.UpdateCurriculum(r.Context(), &c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("curriculum not found for update", slog.Int64("curriculum_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("curriculum not found"))
				return
			}
			log.Error("failed to update curriculum", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update curriculum"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, c)
	}
}

// @Summary Удалить учебный план
// @Tags curriculums
// @Accept json
// @Produce json
// @Param id path int true "ID учебного плана"
// @Success 204 {string} string "No Content"
// @Router /api/v1/curriculums/{id} [delete]
// @Security BearerAuth
func (h *CurriculumHandler) DeleteCurriculum(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.curriculum_handler.DeleteCurriculum"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid curriculum id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid curriculum id"))
			return
		}
		if err := h.repo.DeleteCurriculum(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("curriculum not found for delete", slog.Int64("curriculum_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("curriculum not found"))
				return
			}
			log.Error("failed to delete curriculum", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete curriculum"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список учебных планов
// @Tags curriculums
// @Accept json
// @Produce json
// @Param semester_id query int false "ID семестра"
// @Param discipline_id query int false "ID дисциплины"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Curriculum
// @Router /api/v1/curriculums [get]
// @Security BearerAuth
func (h *CurriculumHandler) ListCurriculum(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.curriculum_handler.ListCurriculum"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var semesterID, disciplineID *int64

		semesterIDStr := r.URL.Query().Get("semester_id")
		if semesterIDStr != "" {
			id, err := strconv.ParseInt(semesterIDStr, 10, 64)
			if err == nil {
				semesterID = &id
			}
		}
		disciplineIDStr := r.URL.Query().Get("discipline_id")
		if disciplineIDStr != "" {
			id, err := strconv.ParseInt(disciplineIDStr, 10, 64)
			if err == nil {
				disciplineID = &id
			}
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 {
			limit = 20
		}

		items, err := h.repo.ListCurriculum(r.Context(), semesterID, disciplineID, limit, offset)
		if err != nil {
			log.Error("failed to list curriculums", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list curriculums"))
			return
		}
		render.JSON(w, r, items)
	}
}
