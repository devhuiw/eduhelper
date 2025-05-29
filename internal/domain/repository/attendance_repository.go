package repository

import (
	"context"
	"service/internal/domain/models"
)

type AttendanceRepository interface {
	Create(ctx context.Context, a *models.Attendance) error
	GetByID(ctx context.Context, id int64) (*models.Attendance, error)
	Update(ctx context.Context, a *models.Attendance) error
	Delete(ctx context.Context, id int64) error
	ListByStudent(ctx context.Context, studentID int64, limit, offset int) ([]*models.Attendance, error)
}
