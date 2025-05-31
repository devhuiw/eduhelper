package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type studentGroupRepository interface {
	CreateStudentGroup(ctx context.Context, group *models.StudentGroup) error
	GetStudentGroupByID(ctx context.Context, id int64) (*models.StudentGroup, error)
	GetStudentGroupPublicByID(ctx context.Context, id int64) (*models.StudentGroupPublic, error)
	UpdateStudentGroup(ctx context.Context, group *models.StudentGroup) error
	DeleteStudentGroup(ctx context.Context, id int64) error
	ListStudentGroups(ctx context.Context, limit, offset int) ([]*models.StudentGroup, error)
	ListStudentGroupPublic(ctx context.Context, limit, offset int) ([]*models.StudentGroupPublic, error)
}

type StudentGroupRepository struct {
	db *sql.DB
}

func NewStudentGroupRepository(db *sql.DB) *StudentGroupRepository {
	return &StudentGroupRepository{db: db}
}

func (r *StudentGroupRepository) CreateStudentGroup(ctx context.Context, group *models.StudentGroup) error {
	query := `
		INSERT INTO student_group (student_group_name, curator_id, academic_year_id, created_at, update_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	group.CreatedAt = now
	group.UpdateAt = now

	result, err := r.db.ExecContext(ctx, query,
		group.StudentGroupName,
		group.CuratorID,
		group.AcademicYearID,
		group.CreatedAt,
		group.UpdateAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		group.StudentGroupID = id
	}
	return err
}

func (r *StudentGroupRepository) GetStudentGroupByID(ctx context.Context, id int64) (*models.StudentGroup, error) {
	query := `
		SELECT student_group_id, created_at, update_at, student_group_name, curator_id, academic_year_id
		FROM student_group
		WHERE student_group_id = ?
	`
	group := &models.StudentGroup{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.StudentGroupID,
		&group.CreatedAt,
		&group.UpdateAt,
		&group.StudentGroupName,
		&group.CuratorID,
		&group.AcademicYearID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return group, nil
}

func (r *StudentGroupRepository) GetStudentGroupPublicByID(ctx context.Context, id int64) (*models.StudentGroupPublic, error) {
	query := `
		SELECT
			sg.student_group_id,
			sg.student_group_name,
			sg.curator_id,
			u.first_name,
			u.last_name,
			u.middle_name,
			sg.academic_year_id
		FROM student_group sg
		JOIN user u ON sg.curator_id = u.user_id
		WHERE sg.student_group_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	group := &models.StudentGroupPublic{}
	var middleName sql.NullString
	err := row.Scan(
		&group.StudentGroupID,
		&group.StudentGroupName,
		&group.CuratorID,
		&group.CuratorFirstName,
		&group.CuratorLastName,
		&middleName,
		&group.AcademicYearID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if middleName.Valid {
		group.CuratorMiddleName = &middleName.String
	}
	return group, nil
}

func (r *StudentGroupRepository) UpdateStudentGroup(ctx context.Context, group *models.StudentGroup) error {
	query := `
		UPDATE student_group
		SET student_group_name = ?, curator_id = ?, academic_year_id = ?, update_at = ?
		WHERE student_group_id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		group.StudentGroupName,
		group.CuratorID,
		group.AcademicYearID,
		time.Now(),
		group.StudentGroupID,
	)
	return err
}

func (r *StudentGroupRepository) DeleteStudentGroup(ctx context.Context, id int64) error {
	query := `DELETE FROM student_group WHERE student_group_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *StudentGroupRepository) ListStudentGroups(ctx context.Context, limit, offset int) ([]*models.StudentGroup, error) {
	query := `
		SELECT student_group_id, created_at, update_at, student_group_name, curator_id, academic_year_id
		FROM student_group
		ORDER BY student_group_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.StudentGroup
	for rows.Next() {
		group := &models.StudentGroup{}
		err := rows.Scan(
			&group.StudentGroupID,
			&group.CreatedAt,
			&group.UpdateAt,
			&group.StudentGroupName,
			&group.CuratorID,
			&group.AcademicYearID,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *StudentGroupRepository) ListStudentGroupPublic(ctx context.Context, limit, offset int) ([]*models.StudentGroupPublic, error) {
	query := `
		SELECT
			sg.student_group_id,
			sg.student_group_name,
			sg.curator_id,
			u.first_name,
			u.last_name,
			u.middle_name,
			sg.academic_year_id
		FROM student_group sg
		JOIN user u ON sg.curator_id = u.user_id
		ORDER BY sg.student_group_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.StudentGroupPublic
	for rows.Next() {
		group := &models.StudentGroupPublic{}
		var middleName sql.NullString
		err := rows.Scan(
			&group.StudentGroupID,
			&group.StudentGroupName,
			&group.CuratorID,
			&group.CuratorFirstName,
			&group.CuratorLastName,
			&middleName,
			&group.AcademicYearID,
		)
		if err != nil {
			return nil, err
		}
		if middleName.Valid {
			group.CuratorMiddleName = &middleName.String
		}
		groups = append(groups, group)
	}
	return groups, nil
}
