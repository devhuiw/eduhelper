package repository

import (
	"context"
	"service/internal/domain/models"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *models.Permission) error
	GetByID(ctx context.Context, id int64) (*models.Permission, error)
	GetByName(ctx context.Context, name string) (*models.Permission, error)
	Update(ctx context.Context, permission *models.Permission) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Permission, error)
}
