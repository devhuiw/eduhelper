package repository

import (
	"context"
	"service/internal/domain/models"
)

type TeacherRepository interface {
	Create(ctx context.Context, teacher *models.Teacher) error
	GetByID(ctx context.Context, userID int64) (*models.Teacher, error)
	Update(ctx context.Context, teacher *models.Teacher) error
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Teacher, error)
}
