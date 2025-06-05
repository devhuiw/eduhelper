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
	"service/internal/lib/utils"
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
	repo      TeacherRepository
	auditRepo AuditLogRepository
}

func NewTeacherHandler(repo TeacherRepository, auditRepo AuditLogRepository) *TeacherHandler {
	return &TeacherHandler{repo: repo, auditRepo: auditRepo}
}

// @Summary Создать преподавателя
// @Tags teachers
// @Accept json
// @Produce json
// @Param input body models.Teacher true "Преподаватель"
// @Success 201 {object} models.Teacher
// @Router /api/v1/teacher [post]
// @Security BearerAuth
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "teacher",
			RowID:      teacher.UserID,
			ActionType: "CREATE",
			NewData:    utils.PtrToJSON(teacher),
			Comment:    utils.PtrToStr("Teacher created"),
		})
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, teacher)
	}
}

// @Summary Получить преподавателя по ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "ID преподавателя"
// @Success 200 {object} models.Teacher
// @Router /api/v1/teacher/{id} [get]
// @Security BearerAuth
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

// @Summary Получить публичный профиль преподавателя по ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "ID преподавателя"
// @Success 200 {object} models.TeacherPublic
// @Router /api/v1/teacher/public/{id} [get]
// @Security BearerAuth
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

// @Summary Получить свой профиль преподавателя
// @Tags teachers
// @Accept json
// @Produce json
// @Success 200 {object} models.Teacher
// @Router /api/v1/teacher/me [get]
// @Security BearerAuth
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

// @Summary Обновить преподавателя по ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "ID преподавателя"
// @Param input body models.Teacher true "Преподаватель"
// @Success 200 {object} models.Teacher
// @Router /api/v1/teacher/{id} [put]
// @Security BearerAuth
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
		oldData, _ := h.repo.GetTeacherByID(r.Context(), teacherId)
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "teacher",
			RowID:      teacherId,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(teacher),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Teacher updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, teacher)
	}
}

// @Summary Обновить свой профиль преподавателя
// @Tags teachers
// @Accept json
// @Produce json
// @Param input body models.Teacher true "Преподаватель"
// @Success 200 {object} models.Teacher
// @Router /api/v1/teacher/me [put]
// @Security BearerAuth
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
		oldData, _ := h.repo.GetTeacherByID(r.Context(), teacherId)
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "teacher",
			RowID:      teacherId,
			ActionType: "UPDATE",
			NewData:    utils.PtrToJSON(teacher),
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Teacher updated"),
		})
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, teacher)
	}
}

// @Summary Удалить преподавателя по ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "ID преподавателя"
// @Success 204 {string} string "No Content"
// @Router /api/v1/teacher/{id} [delete]
// @Security BearerAuth
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
		oldData, _ := h.repo.GetTeacherByID(r.Context(), id)
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
		_ = h.auditRepo.AddAuditLog(r.Context(), &models.AuditLog{
			UserID:     utils.GetUserIDFromContext(r.Context()),
			TableName:  "teacher",
			RowID:      id,
			ActionType: "DELETE",
			OldData:    utils.PtrToJSON(oldData),
			Comment:    utils.PtrToStr("Teacher created"),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Получить список преподавателей
// @Tags teachers
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Teacher
// @Router /api/v1/teacher [get]
// @Security BearerAuth
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

// @Summary Получить публичный список преподавателей
// @Tags teachers
// @Accept json
// @Produce json
// @Param limit query int false "Ограничение"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.TeacherPublic
// @Router /api/v1/teacher/public [get]
// @Security BearerAuth
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
