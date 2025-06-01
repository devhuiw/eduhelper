package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
	"time"
)

type UserRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		userID, roleID, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRoleRepository) RemoveRole(ctx context.Context, userID, roleID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_roles WHERE user_id = ? AND role_id = ?`, userID, roleID)
	return err
}

func (r *UserRoleRepository) GetRolesByUserID(ctx context.Context, userID int64) ([]*models.UserRole, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT created_at, updated_at, role_id, user_id
		 FROM user_roles
		 WHERE user_id = ?`, userID)
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
