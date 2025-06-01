package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) CreateTeacher(ctx context.Context, teacher *models.Teacher) error {
	query := `
		INSERT INTO teacher (user_id, phone, working_experience, education, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	teacher.CreatedAt = now
	teacher.UpdateAt = now

	_, err := r.db.ExecContext(
		ctx, query,
		teacher.UserID,
		teacher.Phone,
		teacher.WorkingExperience,
		teacher.Education,
		teacher.CreatedAt,
		teacher.UpdateAt,
	)
	return err
}

func (r *TeacherRepository) GetTeacherByID(ctx context.Context, userID int64) (*models.Teacher, error) {
	query := `
		SELECT user_id, phone, working_experience, education
		FROM teacher
		WHERE user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	teacher := &models.Teacher{}

	err := row.Scan(
		&teacher.UserID,
		&teacher.Phone,
		&teacher.WorkingExperience,
		&teacher.Education,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return teacher, nil
}

func (r *TeacherRepository) GetTeacherPublicByID(ctx context.Context, userID int64) (*models.TeacherPublic, error) {
	query := `
		SELECT t.user_id, u.first_name, u.last_name, u.middle_name, t.education
		FROM teacher t
		JOIN "user" u ON t.user_id = u.user_id
		WHERE t.user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	teacher := &models.TeacherPublic{}
	var middleName sql.NullString

	err := row.Scan(
		&teacher.UserID,
		&teacher.FirstName,
		&teacher.LastName,
		&middleName,
		&teacher.Education,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if middleName.Valid {
		teacher.MiddleName = &middleName.String
	}
	return teacher, nil
}

func (r *TeacherRepository) UpdateTeacher(ctx context.Context, teacher *models.Teacher) error {
	query := `
		UPDATE teacher SET
			phone = ?, working_experience = ?, education = ?, updated_at = ?
		WHERE user_id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query,
		teacher.Phone,
		teacher.WorkingExperience,
		teacher.Education,
		time.Now(),
		teacher.UserID,
	)
	return err
}

func (r *TeacherRepository) DeleteTeacher(ctx context.Context, userID int64) error {
	query := `DELETE FROM teacher WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *TeacherRepository) ListTeacher(ctx context.Context, limit, offset int) ([]*models.Teacher, error) {
	query := `
		SELECT user_id, phone, working_experience, education
		FROM teacher
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*models.Teacher
	for rows.Next() {
		teacher := &models.Teacher{}
		err := rows.Scan(
			&teacher.UserID,
			&teacher.Phone,
			&teacher.WorkingExperience,
			&teacher.Education,
		)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (r *TeacherRepository) ListTeacherPublic(ctx context.Context, limit, offset int) ([]*models.TeacherPublic, error) {
	query := `
		SELECT t.user_id, u.first_name, u.last_name, u.middle_name, t.education
		FROM teacher t
		INNER JOIN "user" u ON t.user_id = u.user_id
		ORDER BY t.user_id LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*models.TeacherPublic
	for rows.Next() {
		teacher := &models.TeacherPublic{}
		var middleName sql.NullString
		err := rows.Scan(
			&teacher.UserID,
			&teacher.FirstName,
			&teacher.LastName,
			&middleName,
			&teacher.Education,
		)
		if err != nil {
			return nil, err
		}
		if middleName.Valid {
			teacher.MiddleName = &middleName.String
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}
