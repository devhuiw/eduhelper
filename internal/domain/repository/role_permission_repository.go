package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
	"time"
)

type RolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

func (r *RolePermissionRepository) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT (role_id, permission_id) DO NOTHING`,
		roleID, permissionID, time.Now(),
	)
	return err
}

func (r *RolePermissionRepository) RemovePermission(ctx context.Context, roleID, permissionID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`,
		roleID, permissionID,
	)
	return err
}

func (r *RolePermissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.permission_id, p.permission_name, p.created_at, p.updated_at
		 FROM permissions p
		 INNER JOIN role_permissions rp ON rp.permission_id = p.permission_id
		 WHERE rp.role_id = ?`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.PermissionID, &perm.PermissionName, &perm.CreatedAt, &perm.UpdateAt); err != nil {
			return nil, err
		}
		perms = append(perms, &perm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return perms, nil
}
