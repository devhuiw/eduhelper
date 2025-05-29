package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
	"time"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]*models.UserRole, error)
}

type userRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_roles (user_id, role_id, created_at, update_at)
		 VALUES ($1, $2, $3, $3)
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		userID, roleID, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRoleRepository) RemoveRole(ctx context.Context, userID, roleID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`, userID, roleID)
	return err
}

func (r *userRoleRepository) GetRolesByUserID(ctx context.Context, userID int64) ([]*models.UserRole, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT created_at, update_at, role_id, user_id
		 FROM user_roles
		 WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.UserRole
	for rows.Next() {
		var ur models.UserRole
		if err := rows.Scan(&ur.CreatedAt, &ur.UpdateAt, &ur.RoleID, &ur.UserID); err != nil {
			return nil, err
		}
		roles = append(roles, &ur)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
