package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type SemesterRepository interface {
	CreateSemester(ctx context.Context, s *models.Semester) error
	GetSemesterByID(ctx context.Context, id int64) (*models.Semester, error)
	UpdateSemester(ctx context.Context, s *models.Semester) error
	DeleteSemester(ctx context.Context, id int64) error
	ListSemester(ctx context.Context, academicYearID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.Semester, error)
}

type semesterRepository struct {
	db *sql.DB
}

func NewSemesterRepository(db *sql.DB) SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) CreateSemester(ctx context.Context, s *models.Semester) error {
	query := `
		INSERT INTO semester (created_at, updated_at, start_with, ends_with, academic_year_id)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	s.CreatedAt = now
	s.UpdateAt = now
	res, err := r.db.ExecContext(ctx, query, s.CreatedAt, s.UpdateAt, s.StartWith, s.EndsWith, s.AcademicYearID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		s.SemesterID = id
	}
	return err
}

func (r *semesterRepository) GetSemesterByID(ctx context.Context, id int64) (*models.Semester, error) {
	query := `
		SELECT semester_id, created_at, updated_at, start_with, ends_with, academic_year_id
		FROM semester
		WHERE semester_id = ?
	`
	s := &models.Semester{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.SemesterID,
		&s.CreatedAt,
		&s.UpdateAt,
		&s.StartWith,
		&s.EndsWith,
		&s.AcademicYearID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return s, nil
}

func (r *semesterRepository) UpdateSemester(ctx context.Context, s *models.Semester) error {
	query := `
		UPDATE semester
		SET updated_at = ?, start_with = ?, ends_with = ?, academic_year_id = ?
		WHERE semester_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), s.StartWith, s.EndsWith, s.AcademicYearID, s.SemesterID)
	return err
}

func (r *semesterRepository) DeleteSemester(ctx context.Context, id int64) error {
	query := `DELETE FROM semester WHERE semester_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *semesterRepository) ListSemester(
	ctx context.Context,
	academicYearID *int64,
	fromDate, toDate *time.Time,
	limit, offset int,
) ([]*models.Semester, error) {
	query := `SELECT semester_id, created_at, updated_at, start_with, ends_with, academic_year_id FROM semester WHERE 1=1`
	var args []interface{}
	if academicYearID != nil {
		query += " AND academic_year_id = ?"
		args = append(args, *academicYearID)
	}
	if fromDate != nil {
		query += " AND start_with >= ?"
		args = append(args, *fromDate)
	}
	if toDate != nil {
		query += " AND ends_with <= ?"
		args = append(args, *toDate)
	}
	query += " ORDER BY semester_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var semesters []*models.Semester
	for rows.Next() {
		s := &models.Semester{}
		err := rows.Scan(
			&s.SemesterID,
			&s.CreatedAt,
			&s.UpdateAt,
			&s.StartWith,
			&s.EndsWith,
			&s.AcademicYearID,
		)
		if err != nil {
			return nil, err
		}
		semesters = append(semesters, s)
	}
	return semesters, nil
}
