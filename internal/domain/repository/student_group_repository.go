package repository

import (
	"context"
	"service/internal/domain/models"
)

type StudentGroupRepository interface {
	Create(ctx context.Context, group *models.StudentGroup) error
	GetByID(ctx context.Context, id int64) (*models.StudentGroup, error)
	Update(ctx context.Context, group *models.StudentGroup) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.StudentGroup, error)
}
