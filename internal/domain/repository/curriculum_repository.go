package repository

import (
	"context"
	"service/internal/domain/models"
)

type CurriculumRepository interface {
	Create(ctx context.Context, c *models.Curriculum) error
	GetByID(ctx context.Context, id int64) (*models.Curriculum, error)
	Update(ctx context.Context, c *models.Curriculum) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Curriculum, error)
}
