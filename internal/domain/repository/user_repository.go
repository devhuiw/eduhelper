package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateClient(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO user (
			first_name, last_name, middle_name, email, password, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdateAt = now

	res, err := r.db.ExecContext(
		ctx, query,
		user.FirstName,
		user.LastName,
		user.MiddleName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdateAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.UserID = id
	return nil
}

func (r *UserRepository) GetClientByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT user_id, created_at, updated_at, first_name, last_name, middle_name, email, password
		FROM user WHERE user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	var middleName sql.NullString

	err := row.Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.UpdateAt,
		&user.FirstName,
		&user.LastName,
		&middleName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if middleName.Valid {
		user.MiddleName = &middleName.String
	}
	return user, nil
}

func (r *UserRepository) GetClientByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT user_id, created_at, updated_at, first_name, last_name, middle_name, email, password
		FROM user WHERE email = ?
	`
	row := r.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	var middleName sql.NullString

	err := row.Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.UpdateAt,
		&user.FirstName,
		&user.LastName,
		&middleName,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if middleName.Valid {
		user.MiddleName = &middleName.String
	}
	return user, nil
}

func (r *UserRepository) UpdateClient(ctx context.Context, user *models.User) error {
	query := `
		UPDATE user SET
			first_name = ?, last_name = ?, middle_name = ?, email = ?, password = ?, updated_at = ?
		WHERE user_id = ?
	`
	user.UpdateAt = time.Now()
	_, err := r.db.ExecContext(
		ctx, query,
		user.FirstName,
		user.LastName,
		user.MiddleName,
		user.Email,
		user.Password,
		user.UpdateAt,
		user.UserID,
	)
	return err
}

func (r *UserRepository) DeleteClient(ctx context.Context, id int64) error {
	query := `DELETE FROM user WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) ListClient(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT user_id, created_at, updated_at, first_name, last_name, middle_name, email, password
		FROM user ORDER BY user_id LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var middleName sql.NullString
		err := rows.Scan(
			&user.UserID,
			&user.CreatedAt,
			&user.UpdateAt,
			&user.FirstName,
			&user.LastName,
			&middleName,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, err
		}
		if middleName.Valid {
			user.MiddleName = &middleName.String
		}
		users = append(users, user)
	}
	return users, nil
}
