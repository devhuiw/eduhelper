package repository

import (
	"context"
	"service/internal/domain/models"
)

type SemesterRepository interface {
	Create(ctx context.Context, semester *models.Semester) error
	GetByID(ctx context.Context, id int64) (*models.Semester, error)
	Update(ctx context.Context, semester *models.Semester) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Semester, error)
}
