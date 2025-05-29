package repository

import (
	"context"
	"service/internal/domain/models"
)

type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, userID int64) (*models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Student, error)
}
