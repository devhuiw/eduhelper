package repository

import (
	"context"
	"service/internal/domain/models"
)

type RolePermissionRepository interface {
	AssignPermission(ctx context.Context, roleID, permissionID int64) error
	RemovePermission(ctx context.Context, roleID, permissionID int64) error
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]*models.Permission, error)
}
