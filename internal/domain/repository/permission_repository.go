package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
	"time"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	query := `
		INSERT INTO permissions (permission_name, created_at, updated_at)
		VALUES (?, ?, ?)
		RETURNING permission_id
	`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, permission.PermissionName, now).Scan(&permission.PermissionID)
	return err
}

func (r *PermissionRepository) GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error) {
	query := `
		SELECT permission_id, permission_name, created_at, updated_at
		FROM permissions
		WHERE permission_id = ?
	`
	var perm models.Permission
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&perm.PermissionID,
		&perm.PermissionName,
		&perm.CreatedAt,
		&perm.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *PermissionRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	query := `
		SELECT permission_id, permission_name, created_at, updated_at
		FROM permissions
		WHERE permission_name = ?
	`
	var perm models.Permission
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&perm.PermissionID,
		&perm.PermissionName,
		&perm.CreatedAt,
		&perm.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *PermissionRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	query := `
		UPDATE permissions
		SET permission_name = ?, updated_at = ?
		WHERE permission_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, permission.PermissionName, time.Now(), permission.PermissionID)
	return err
}

func (r *PermissionRepository) DeletePermission(ctx context.Context, id int64) error {
	query := `DELETE FROM permissions WHERE permission_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PermissionRepository) ListPermission(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	query := `
		SELECT permission_id, permission_name, created_at, updated_at
		FROM permissions
		ORDER BY permission_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
