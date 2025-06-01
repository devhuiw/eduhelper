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

type AcademicYearRepository interface {
	CreateAcademicYear(ctx context.Context, year *models.AcademicYear) error
	GetAcademicYearByID(ctx context.Context, id int64) (*models.AcademicYear, error)
	UpdateAcademicYear(ctx context.Context, year *models.AcademicYear) error
	DeleteAcademicYear(ctx context.Context, id int64) error
	ListAcademicYear(ctx context.Context, limit, offset int) ([]*models.AcademicYear, error)
}

type AcademicYearHandler struct {
	repo AcademicYearRepository
}

func NewAcademicYearHandler(repo AcademicYearRepository) *AcademicYearHandler {
	return &AcademicYearHandler{repo: repo}
}

// @Summary Создать учебный год
// @Tags academic-years
// @Accept json
// @Produce json
// @Param input body models.AcademicYear true "Учебный год"
// @Success 201 {object} models.AcademicYear
// @Router /api/v1/academic-years [post]
// @Security BearerAuth
func (h *AcademicYearHandler) CreateAcademicYear(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.academicyear_handler.CreateAcademicYear"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var year models.AcademicYear
		if err := json.NewDecoder(r.Body).Decode(&year); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := h.repo.CreateAcademicYear(r.Context(), &year); err != nil {
			log.Error("failed to create academic year", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create academic year"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, year)
	}
}

// @Summary Получить учебный год по ID
// @Tags academic-years
// @Accept json
// @Produce json
// @Param id path int true "ID учебного года"
// @Success 200 {object} models.AcademicYear
// @Router /api/v1/academic-years/{id} [get]
// @Security BearerAuth
func (h *AcademicYearHandler) GetAcademicYearByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.academicyear_handler.GetAcademicYearByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid academic year id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid academic year id"))
			return
		}
		year, err := h.repo.GetAcademicYearByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("academic year not found", slog.Int64("academic_year_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("academic year not found"))
				return
			}
			log.Error("failed to get academic year", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get academic year"))
			return
		}
		render.JSON(w, r, year)
	}
}

// @Summary Обновить учебный год
// @Tags academic-years
// @Accept json
// @Produce json
// @Param id path int true "ID учебного года"
// @Param input body models.AcademicYear true "Учебный год"
// @Success 200 {object} models.AcademicYear
// @Router /api/v1/academic-years/{id} [put]
// @Security BearerAuth
func (h *AcademicYearHandler) UpdateAcademicYear(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.academicyear_handler.UpdateAcademicYear"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid academic year id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid academic year id"))
			return
		}
		var year models.AcademicYear
		if err := json.NewDecoder(r.Body).Decode(&year); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		year.AcademicYearID = id
		if err := h.repo.UpdateAcademicYear(r.Context(), &year); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("academic year not found for update", slog.Int64("academic_year_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("academic year not found"))
				return
			}
			log.Error("failed to update academic year", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update academic year"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, year)
	}
}

// @Summary Удалить учебный год
// @Tags academic-years
// @Accept json
// @Produce json
// @Param id path int true "ID учебного года"
// @Success 204 {string} string "No Content"
// @Router /api/v1/academic-years/{id} [delete]
// @Security BearerAuth
func (h *AcademicYearHandler) DeleteAcademicYear(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.academicyear_handler.DeleteAcademicYear"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid academic year id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid academic year id"))
			return
		}
		if err := h.repo.DeleteAcademicYear(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("academic year not found for delete", slog.Int64("academic_year_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("academic year not found"))
				return
			}
			log.Error("failed to delete academic year", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete academic year"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список учебных годов
// @Tags academic-years
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.AcademicYear
// @Router /api/v1/academic-years [get]
// @Security BearerAuth
func (h *AcademicYearHandler) ListAcademicYear(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.academicyear_handler.ListAcademicYear"
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
		years, err := h.repo.ListAcademicYear(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list academic years", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list academic years"))
			return
		}
		render.JSON(w, r, years)
	}
}
