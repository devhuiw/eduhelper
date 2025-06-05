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

type StudentGroupRepository interface {
	CreateStudentGroup(ctx context.Context, group *models.StudentGroup) error
	GetStudentGroupByID(ctx context.Context, id int64) (*models.StudentGroup, error)
	GetStudentGroupPublicByID(ctx context.Context, id int64) (*models.StudentGroupPublic, error)
	UpdateStudentGroup(ctx context.Context, group *models.StudentGroup) error
	DeleteStudentGroup(ctx context.Context, id int64) error
	ListStudentGroups(ctx context.Context, limit, offset int) ([]*models.StudentGroup, error)
	ListStudentGroupPublic(ctx context.Context, limit, offset int) ([]*models.StudentGroupPublic, error)
}

type StudentGroupHandler struct {
	repo      StudentGroupRepository
	auditRepo AuditLogRepository
}

func NewStudentGroupHandler(repo StudentGroupRepository, auditRepo AuditLogRepository) *StudentGroupHandler {
	return &StudentGroupHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать группу студентов
// @Tags student-groups
// @Accept json
// @Produce json
// @Param input body models.StudentGroup true "Группа"
// @Success 201 {object} models.StudentGroup
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/v1/student-groups [post]
// @Security BearerAuth
func (h *StudentGroupHandler) CreateStudentGroup(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.CreateStudentGroup"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var group models.StudentGroup
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		if err := h.repo.CreateStudentGroup(r.Context(), &group); err != nil {
			log.Error("failed to create student group", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create student group"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student_group",
			RowID:      group.StudentGroupID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(group),
			Comment:    utils.PtrToStr("Student group created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, group)
	}
}

// @Summary Получить группу по ID
// @Tags student-groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Success 200 {object} models.StudentGroup
// @Failure 404 {object} resp.Response
// @Router /api/v1/student-groups/{id} [get]
// @Security BearerAuth
func (h *StudentGroupHandler) GetStudentGroupByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.GetStudentGroupByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid student group id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid group id"))
			return
		}
		group, err := h.repo.GetStudentGroupByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student group not found", slog.Int64("student_group_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("group not found"))
				return
			}
			log.Error("failed to get student group", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get group"))
			return
		}
		render.JSON(w, r, group)
	}
}

// @Summary Получить публичную группу по ID
// @Tags student-groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Success 200 {object} models.StudentGroupPublic
// @Failure 404 {object} resp.Response
// @Router /api/v1/student-groups/public/{id} [get]
// @Security BearerAuth
func (h *StudentGroupHandler) GetStudentGroupPublicByID(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.GetStudentGroupPublicByID"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid group id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid group id"))
			return
		}
		group, err := h.repo.GetStudentGroupPublicByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("student group not found", slog.Int64("student_group_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("group not found"))
				return
			}
			log.Error("failed to get group public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get group"))
			return
		}
		render.JSON(w, r, group)
	}
}

// @Summary Обновить группу студентов
// @Tags student-groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Param input body models.StudentGroup true "Группа"
// @Success 200 {object} models.StudentGroup
// @Failure 404 {object} resp.Response
// @Router /api/v1/student-groups/{id} [put]
// @Security BearerAuth
func (h *StudentGroupHandler) UpdateStudentGroup(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.UpdateStudentGroup"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid group id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid group id"))
			return
		}
		var group models.StudentGroup
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			log.Info("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		group.StudentGroupID = id
		oldData, _ := h.repo.GetStudentGroupByID(r.Context(), id)
		if err := h.repo.UpdateStudentGroup(r.Context(), &group); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("group not found for update", slog.Int64("student_group_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("group not found"))
				return
			}
			log.Error("failed to update group", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update group"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student_group",
			RowID:      id,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(group),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Student Group updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, group)
	}
}

// @Summary Удалить группу студентов
// @Tags student-groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} resp.Response
// @Router /api/v1/student-groups/{id} [delete]
// @Security BearerAuth
func (h *StudentGroupHandler) DeleteStudentGroup(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.DeleteStudentGroup"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid group id", slog.String("id", idStr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid group id"))
			return
		}
		oldData, _ := h.repo.GetStudentGroupByID(r.Context(), id)
		if err := h.repo.DeleteStudentGroup(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("group not found for delete", slog.Int64("student_group_id", id))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("group not found"))
				return
			}
			log.Error("failed to delete group", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete group"))
			return
		}
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "student_group",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Student Group deleted"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список групп студентов
// @Tags student-groups
// @Accept json
// @Produce json
// @Param curator_id query int false "ID куратора"
// @Param academic_year_id query int false "ID учебного года"
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.StudentGroup
// @Router /api/v1/student-groups [get]
// @Security BearerAuth
func (h *StudentGroupHandler) ListStudentGroups(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.ListStudentGroups"
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
		groups, err := h.repo.ListStudentGroups(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list groups", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list groups"))
			return
		}
		render.JSON(w, r, groups)
	}
}

// @Summary Получить список публичных групп
// @Tags student-groups
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.StudentGroupPublic
// @Router /api/v1/student-groups/public [get]
// @Security BearerAuth
func (h *StudentGroupHandler) ListStudentGroupPublic(log *slog.Logger) http.HandlerFunc {
	const op = "handler.v1.studentgroup_handler.ListStudentGroupPublic"
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
		groups, err := h.repo.ListStudentGroupPublic(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to list groups public", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to list groups public"))
			return
		}
		render.JSON(w, r, groups)
	}
}
