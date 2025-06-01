package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type CurriculumRepository interface {
	CreateCurriculum(ctx context.Context, c *models.Curriculum) error
	GetCurriculumByID(ctx context.Context, id int64) (*models.Curriculum, error)
	UpdateCurriculum(ctx context.Context, c *models.Curriculum) error
	DeleteCurriculum(ctx context.Context, id int64) error
	ListCurriculum(ctx context.Context, semesterID, disciplineID *int64, limit, offset int) ([]*models.Curriculum, error)
}

type curriculumRepository struct {
	db *sql.DB
}

func NewCurriculumRepository(db *sql.DB) CurriculumRepository {
	return &curriculumRepository{db: db}
}

func (r *curriculumRepository) CreateCurriculum(ctx context.Context, c *models.Curriculum) error {
	query := `
		INSERT INTO curriculum (created_at, updated_at, subject_name, subject_description, semester_id, discipline_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	c.CreatedAt = now
	c.UpdateAt = now
	res, err := r.db.ExecContext(ctx, query, c.CreatedAt, c.UpdateAt, c.SubjectName, c.SubjectDescription, c.SemesterID, c.DisciplineID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		c.CurriculumID = id
	}
	return err
}

func (r *curriculumRepository) GetCurriculumByID(ctx context.Context, id int64) (*models.Curriculum, error) {
	query := `
		SELECT curriculum_id, created_at, updated_at, subject_name, subject_description, semester_id, discipline_id
		FROM curriculum WHERE curriculum_id = ?
	`
	c := &models.Curriculum{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.CurriculumID,
		&c.CreatedAt,
		&c.UpdateAt,
		&c.SubjectName,
		&c.SubjectDescription,
		&c.SemesterID,
		&c.DisciplineID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return c, nil
}

func (r *curriculumRepository) UpdateCurriculum(ctx context.Context, c *models.Curriculum) error {
	query := `
		UPDATE curriculum
		SET updated_at = ?, subject_name = ?, subject_description = ?, semester_id = ?, discipline_id = ?
		WHERE curriculum_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), c.SubjectName, c.SubjectDescription, c.SemesterID, c.DisciplineID, c.CurriculumID)
	return err
}

func (r *curriculumRepository) DeleteCurriculum(ctx context.Context, id int64) error {
	query := `DELETE FROM curriculum WHERE curriculum_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *curriculumRepository) ListCurriculum(
	ctx context.Context,
	semesterID, disciplineID *int64,
	limit, offset int,
) ([]*models.Curriculum, error) {
	query := `SELECT curriculum_id, created_at, updated_at, subject_name, subject_description, semester_id, discipline_id FROM curriculum WHERE 1=1`
	var args []interface{}
	if semesterID != nil {
		query += " AND semester_id = ?"
		args = append(args, *semesterID)
	}
	if disciplineID != nil {
		query += " AND discipline_id = ?"
		args = append(args, *disciplineID)
	}
	query += " ORDER BY curriculum_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Curriculum
	for rows.Next() {
		c := &models.Curriculum{}
		err := rows.Scan(
			&c.CurriculumID,
			&c.CreatedAt,
			&c.UpdateAt,
			&c.SubjectName,
			&c.SubjectDescription,
			&c.SemesterID,
			&c.DisciplineID,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}
