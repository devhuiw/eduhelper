package repository

import (
	"context"
	"service/internal/domain/models"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id int64) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Role, error)
}
