package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type academicYearRepository struct {
	db *sql.DB
}

func NewAcademicYearRepository(db *sql.DB) *academicYearRepository {
	return &academicYearRepository{db: db}
}

func (r *academicYearRepository) CreateAcademicYear(ctx context.Context, year *models.AcademicYear) error {
	query := `
		INSERT INTO academic_year (name_academic_year, start_with, ends_with, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	year.CreatedAt = now
	year.UpdateAt = now

	res, err := r.db.ExecContext(ctx, query,
		year.Name,
		year.StartWith,
		year.EndsWith,
		year.CreatedAt,
		year.UpdateAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		year.AcademicYearID = id
	}
	return err
}

func (r *academicYearRepository) GetAcademicYearByID(ctx context.Context, id int64) (*models.AcademicYear, error) {
	query := `
		SELECT academic_year_id, name_academic_year, start_with, ends_with, created_at, updated_at
		FROM academic_year
		WHERE academic_year_id = ?
	`
	year := &models.AcademicYear{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&year.AcademicYearID,
		&year.Name,
		&year.StartWith,
		&year.EndsWith,
		&year.CreatedAt,
		&year.UpdateAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return year, nil
}

func (r *academicYearRepository) UpdateAcademicYear(ctx context.Context, year *models.AcademicYear) error {
	query := `
		UPDATE academic_year
		SET name_academic_year = ?, start_with = ?, ends_with = ?, updated_at = ?
		WHERE academic_year_id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		year.Name,
		year.StartWith,
		year.EndsWith,
		time.Now(),
		year.AcademicYearID,
	)
	return err
}

func (r *academicYearRepository) DeleteAcademicYear(ctx context.Context, id int64) error {
	query := `DELETE FROM academic_year WHERE academic_year_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *academicYearRepository) ListAcademicYear(ctx context.Context, limit, offset int) ([]*models.AcademicYear, error) {
	query := `
		SELECT academic_year_id, name_academic_year, start_with, ends_with, created_at, updated_at
		FROM academic_year
		ORDER BY academic_year_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []*models.AcademicYear
	for rows.Next() {
		year := &models.AcademicYear{}
		err := rows.Scan(
			&year.AcademicYearID,
			&year.Name,
			&year.StartWith,
			&year.EndsWith,
			&year.CreatedAt,
			&year.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	return years, nil
}
