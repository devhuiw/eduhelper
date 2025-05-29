package repository

import (
	"context"
	"service/internal/domain/models"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]*models.Role, error)
}
