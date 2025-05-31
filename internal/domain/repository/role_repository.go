package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
	"time"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *models.Role) (int64, error) {
	query := `
		INSERT INTO roles (role_name, created_at, update_at)
		VALUES ($1, $2, $3)
		RETURNING role_id
	`
	var id int64
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, role.RoleName, now, now).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id int64) (*models.Role, error) {
	query := `
		SELECT role_id, role_name, created_at, update_at
		FROM roles
		WHERE role_id = $1
	`
	var role models.Role
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.RoleID,
		&role.RoleName,
		&role.CreatedAt,
		&role.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT role_id, role_name, created_at, update_at
		FROM roles
		WHERE role_name = $1
	`
	var role models.Role
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.RoleID,
		&role.RoleName,
		&role.CreatedAt,
		&role.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	query := `
		UPDATE roles
		SET role_name = $1, update_at = $2
		WHERE role_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, role.RoleName, time.Now(), role.RoleID)
	return err
}

func (r *RoleRepository) DeleteRole(ctx context.Context, id int64) error {
	query := `DELETE FROM roles WHERE role_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RoleRepository) ListRole(ctx context.Context) ([]*models.Role, error) {
	query := `
		SELECT role_id, role_name, created_at, update_at
		FROM roles
		ORDER BY role_id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.RoleID, &role.RoleName, &role.CreatedAt, &role.UpdateAt); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, rows.Err()
}
