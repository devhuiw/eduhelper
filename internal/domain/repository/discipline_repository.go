package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"strings"
	"time"
)

type disciplineRepository struct {
	db *sql.DB
}

func NewDisciplineRepository(db *sql.DB) *disciplineRepository {
	return &disciplineRepository{db: db}
}

func (r *disciplineRepository) CreateDiscipline(ctx context.Context, d *models.Discipline) error {
	query := `
		INSERT INTO discipline (discipline_name, teacher_id, student_group_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	d.CreatedAt = now
	d.UpdateAt = now

	res, err := r.db.ExecContext(ctx, query, d.DisciplineName, d.TeacherID, d.StudentGroupID, d.CreatedAt, d.UpdateAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		d.DisciplineID = id
	}
	return err
}

func (r *disciplineRepository) GetDisciplineByID(ctx context.Context, id int64) (*models.Discipline, error) {
	query := `
		SELECT discipline_id, created_at, updated_at, discipline_name, teacher_id, student_group_id
		FROM discipline
		WHERE discipline_id = ?
	`
	d := &models.Discipline{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.DisciplineID,
		&d.CreatedAt,
		&d.UpdateAt,
		&d.DisciplineName,
		&d.TeacherID,
		&d.StudentGroupID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return d, nil
}

func (r *disciplineRepository) UpdateDiscipline(ctx context.Context, d *models.Discipline) error {
	query := `
		UPDATE discipline
		SET discipline_name = ?, teacher_id = ?, student_group_id = ?, updated_at = ?
		WHERE discipline_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, d.DisciplineName, d.TeacherID, d.StudentGroupID, time.Now(), d.DisciplineID)
	return err
}

func (r *disciplineRepository) DeleteDiscipline(ctx context.Context, id int64) error {
	query := `DELETE FROM discipline WHERE discipline_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *disciplineRepository) ListDiscipline(ctx context.Context, limit, offset int) ([]*models.Discipline, error) {
	query := `
		SELECT discipline_id, created_at, updated_at, discipline_name, teacher_id, student_group_id
		FROM discipline
		ORDER BY discipline_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disciplines []*models.Discipline
	for rows.Next() {
		d := &models.Discipline{}
		err := rows.Scan(
			&d.DisciplineID,
			&d.CreatedAt,
			&d.UpdateAt,
			&d.DisciplineName,
			&d.TeacherID,
			&d.StudentGroupID,
		)
		if err != nil {
			return nil, err
		}
		disciplines = append(disciplines, d)
	}
	return disciplines, nil
}

// --- PUBLIC ---

func (r *disciplineRepository) GetDisciplinePublicByID(ctx context.Context, id int64) (*models.DisciplinePublic, error) {
	query := `
SELECT
    d.discipline_id,
    d.created_at,
    d.updated_at,
    d.discipline_name,
    d.teacher_id,
    t.first_name,
    t.last_name,
    t.middle_name,
    d.student_group_id,
    sg.student_group_name,
    sg.curator_id,
    c.first_name AS curator_first_name,
    c.last_name AS curator_last_name,
    c.middle_name AS curator_middle_name,
    sg.academic_year_id
FROM discipline d
JOIN user t ON d.teacher_id = t.user_id
JOIN student_group sg ON d.student_group_id = sg.student_group_id
JOIN user c ON sg.curator_id = c.user_id
WHERE d.discipline_id = ?
`
	dp := &models.DisciplinePublic{}
	var teacherMiddle, curatorMiddle sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dp.DisciplineID,
		&dp.CreatedAt,
		&dp.UpdateAt,
		&dp.DisciplineName,
		&dp.TeacherID,
		&dp.FirstName,
		&dp.LastName,
		&teacherMiddle,
		&dp.StudentGroupID,
		&dp.StudentGroupName,
		&dp.CuratorID,
		&dp.CuratorFirstName,
		&dp.CuratorLastName,
		&curatorMiddle,
		&dp.AcademicYearID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if teacherMiddle.Valid {
		dp.MiddleName = &teacherMiddle.String
	}
	if curatorMiddle.Valid {
		dp.CuratorMiddleName = &curatorMiddle.String
	}
	return dp, nil
}

func (r *disciplineRepository) ListDisciplinePublic(
	ctx context.Context,
	limit, offset int,
	teacherID, studentGroupID, academicYearID *int64,
) ([]*models.DisciplinePublic, error) {
	query := `
		SELECT
			d.discipline_id,
			d.created_at,
			d.updated_at,
			d.discipline_name,
			d.teacher_id,
			t.first_name,
			t.last_name,
			t.middle_name,
			d.student_group_id,
			sg.student_group_name,
			sg.curator_id,
			c.first_name AS curator_first_name,
			c.last_name AS curator_last_name,
			c.middle_name AS curator_middle_name,
			sg.academic_year_id
		FROM discipline d
		JOIN user t ON d.teacher_id = t.user_id
		JOIN student_group sg ON d.student_group_id = sg.student_group_id
		JOIN user c ON sg.curator_id = c.user_id
		`
	var (
		where []string
		args  []interface{}
	)

	if teacherID != nil {
		where = append(where, "d.teacher_id = ?")
		args = append(args, *teacherID)
	}
	if studentGroupID != nil {
		where = append(where, "d.student_group_id = ?")
		args = append(args, *studentGroupID)
	}
	if academicYearID != nil {
		where = append(where, "sg.academic_year_id = ?")
		args = append(args, *academicYearID)
	}

	if len(where) > 0 {
		query += " WHERE " + joinWithAnd(where)
	}
	query += " ORDER BY d.discipline_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disciplines []*models.DisciplinePublic
	for rows.Next() {
		dp := &models.DisciplinePublic{}
		var teacherMiddle, curatorMiddle sql.NullString
		err := rows.Scan(
			&dp.DisciplineID,
			&dp.CreatedAt,
			&dp.UpdateAt,
			&dp.DisciplineName,
			&dp.TeacherID,
			&dp.FirstName,
			&dp.LastName,
			&teacherMiddle,
			&dp.StudentGroupID,
			&dp.StudentGroupName,
			&dp.CuratorID,
			&dp.CuratorFirstName,
			&dp.CuratorLastName,
			&curatorMiddle,
			&dp.AcademicYearID,
		)
		if err != nil {
			return nil, err
		}
		if teacherMiddle.Valid {
			dp.MiddleName = &teacherMiddle.String
		}
		if curatorMiddle.Valid {
			dp.CuratorMiddleName = &curatorMiddle.String
		}
		disciplines = append(disciplines, dp)
	}
	return disciplines, nil
}

func joinWithAnd(conds []string) string {
	return strings.Join(conds, " AND ")
}
