package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, student *models.Student) error {
	query := `
		INSERT INTO student (user_id, phone, birthday, created_at, updated_at, student_group_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	student.CreatedAt = now
	student.UpdateAt = now

	_, err := r.db.ExecContext(
		ctx, query,
		student.UserID,
		student.Phone,
		student.Birthday,
		student.CreatedAt,
		student.UpdateAt,
		student.StudentGroupID,
	)
	return err
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, userID int64) (*models.Student, error) {
	query := `
		SELECT user_id, phone, birthday, created_at, updated_at, student_group_id
		FROM student
		WHERE user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	student := &models.Student{}

	err := row.Scan(
		&student.UserID,
		&student.Phone,
		&student.Birthday,
		&student.CreatedAt,
		&student.UpdateAt,
		&student.StudentGroupID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return student, nil
}

func (r *StudentRepository) GetStudentPublicByID(ctx context.Context, userID int64) (*models.StudentPublic, error) {
	query := `
		SELECT s.user_id, u.first_name, u.last_name, u.middle_name, s.birthday, s.student_group_id
		FROM student s
		JOIN user u ON s.user_id = u.user_id
		WHERE s.user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	student := &models.StudentPublic{}
	var middleName sql.NullString

	err := row.Scan(
		&student.UserID,
		&student.FirstName,
		&student.LastName,
		&middleName,
		&student.Birthday,
		&student.StudentGroupID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if middleName.Valid {
		student.MiddleName = &middleName.String
	}
	return student, nil
}

func (r *StudentRepository) UpdateStudent(ctx context.Context, student *models.Student) error {
	query := `
		UPDATE student SET
			phone = ?, birthday = ?, updated_at = ?, student_group_id = ?
		WHERE user_id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query,
		student.Phone,
		student.Birthday,
		time.Now(),
		student.StudentGroupID,
		student.UserID,
	)
	return err
}

func (r *StudentRepository) DeleteStudent(ctx context.Context, userID int64) error {
	query := `DELETE FROM student WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *StudentRepository) ListStudent(ctx context.Context, limit, offset int) ([]*models.Student, error) {
	query := `
		SELECT user_id, phone, birthday, created_at, updated_at, student_group_id
		FROM student
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		student := &models.Student{}
		err := rows.Scan(
			&student.UserID,
			&student.Phone,
			&student.Birthday,
			&student.CreatedAt,
			&student.UpdateAt,
			&student.StudentGroupID,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (r *StudentRepository) ListStudentPublic(ctx context.Context, limit, offset int) ([]*models.StudentPublic, error) {
	query := `
		SELECT s.user_id, u.first_name, u.last_name, u.middle_name, s.birthday, s.student_group_id
		FROM student s
		INNER JOIN user u ON s.user_id = u.user_id
		ORDER BY s.user_id LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.StudentPublic
	for rows.Next() {
		student := &models.StudentPublic{}
		var middleName sql.NullString
		err := rows.Scan(
			&student.UserID,
			&student.FirstName,
			&student.LastName,
			&middleName,
			&student.Birthday,
			&student.StudentGroupID,
		)
		if err != nil {
			return nil, err
		}
		if middleName.Valid {
			student.MiddleName = &middleName.String
		}
		students = append(students, student)
	}
	return students, nil
}
