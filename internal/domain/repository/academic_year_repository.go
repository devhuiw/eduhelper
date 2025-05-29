package repository

import (
	"context"
	"service/internal/domain/models"
)

type AcademicYearRepository interface {
	Create(ctx context.Context, year *models.AcademicYear) error
	GetByID(ctx context.Context, id int64) (*models.AcademicYear, error)
	Update(ctx context.Context, year *models.AcademicYear) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.AcademicYear, error)
}
