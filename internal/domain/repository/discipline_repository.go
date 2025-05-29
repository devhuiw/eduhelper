package repository

import (
	"context"
	"service/internal/domain/models"
)

type DisciplineRepository interface {
	Create(ctx context.Context, d *models.Discipline) error
	GetByID(ctx context.Context, id int64) (*models.Discipline, error)
	Update(ctx context.Context, d *models.Discipline) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Discipline, error)
}
