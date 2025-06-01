package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type GradeJournalRepository interface {
	CreateGradeJournal(ctx context.Context, g *models.GradeJournal) error
	GetGradeJournalByID(ctx context.Context, id int64) (*models.GradeJournal, error)
	UpdateGradeJournal(ctx context.Context, g *models.GradeJournal) error
	DeleteGradeJournal(ctx context.Context, id int64) error
	ListGradeJournal(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.GradeJournal, error)
	ListGradeJournalPublic(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time, limit, offset int) ([]*models.GradeJournalPublic, error)
	GetAverageGrade(ctx context.Context, studentID, disciplineID *int64, fromDate, toDate *time.Time) (float64, error)
}

type gradeJournalRepository struct {
	db *sql.DB
}

func NewGradeJournalRepository(db *sql.DB) GradeJournalRepository {
	return &gradeJournalRepository{db: db}
}

func (r *gradeJournalRepository) CreateGradeJournal(ctx context.Context, g *models.GradeJournal) error {
	query := `
		INSERT INTO grade_journal (created_at, updated_at, student_id, grade, comment, discipline_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	g.CreatedAt = now
	g.UpdateAt = now
	res, err := r.db.ExecContext(ctx, query, g.CreatedAt, g.UpdateAt, g.StudentID, g.Grade, g.Comment, g.DisciplineID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		g.GradeJournalID = id
	}
	return err
}

func (r *gradeJournalRepository) GetGradeJournalByID(ctx context.Context, id int64) (*models.GradeJournal, error) {
	query := `
		SELECT grade_journal_id, created_at, updated_at, student_id, grade, comment, discipline_id
		FROM grade_journal WHERE grade_journal_id = ?
	`
	g := &models.GradeJournal{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&g.GradeJournalID, &g.CreatedAt, &g.UpdateAt, &g.StudentID, &g.Grade, &g.Comment, &g.DisciplineID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return g, nil
}

func (r *gradeJournalRepository) UpdateGradeJournal(ctx context.Context, g *models.GradeJournal) error {
	query := `
		UPDATE grade_journal SET updated_at = ?, student_id = ?, grade = ?, comment = ?, discipline_id = ?
		WHERE grade_journal_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), g.StudentID, g.Grade, g.Comment, g.DisciplineID, g.GradeJournalID)
	return err
}

func (r *gradeJournalRepository) DeleteGradeJournal(ctx context.Context, id int64) error {
	query := `DELETE FROM grade_journal WHERE grade_journal_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *gradeJournalRepository) ListGradeJournal(
	ctx context.Context,
	studentID, disciplineID *int64,
	fromDate, toDate *time.Time,
	limit, offset int,
) ([]*models.GradeJournal, error) {
	query := `SELECT grade_journal_id, created_at, updated_at, student_id, grade, comment, discipline_id FROM grade_journal WHERE 1=1`
	var args []interface{}
	if studentID != nil {
		query += " AND student_id = ?"
		args = append(args, *studentID)
	}
	if disciplineID != nil {
		query += " AND discipline_id = ?"
		args = append(args, *disciplineID)
	}
	if fromDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *fromDate)
	}
	if toDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *toDate)
	}
	query += " ORDER BY grade_journal_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.GradeJournal
	for rows.Next() {
		g := &models.GradeJournal{}
		err := rows.Scan(
			&g.GradeJournalID,
			&g.CreatedAt,
			&g.UpdateAt,
			&g.StudentID,
			&g.Grade,
			&g.Comment,
			&g.DisciplineID,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, g)
	}
	return items, nil
}

// Публичная версия — join к user и discipline
func (r *gradeJournalRepository) ListGradeJournalPublic(
	ctx context.Context,
	studentID, disciplineID *int64,
	fromDate, toDate *time.Time,
	limit, offset int,
) ([]*models.GradeJournalPublic, error) {
	query := `
		SELECT 
			gj.grade_journal_id, gj.created_at, gj.updated_at, gj.student_id,
			u.first_name, u.last_name,
			gj.discipline_id, d.discipline_name,
			gj.grade, gj.comment
		FROM grade_journal gj
		JOIN user u ON gj.student_id = u.user_id
		JOIN discipline d ON gj.discipline_id = d.discipline_id
		WHERE 1=1
	`
	var args []interface{}
	if studentID != nil {
		query += " AND gj.student_id = ?"
		args = append(args, *studentID)
	}
	if disciplineID != nil {
		query += " AND gj.discipline_id = ?"
		args = append(args, *disciplineID)
	}
	if fromDate != nil {
		query += " AND gj.created_at >= ?"
		args = append(args, *fromDate)
	}
	if toDate != nil {
		query += " AND gj.created_at <= ?"
		args = append(args, *toDate)
	}
	query += " ORDER BY gj.grade_journal_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.GradeJournalPublic
	for rows.Next() {
		g := &models.GradeJournalPublic{}
		err := rows.Scan(
			&g.GradeJournalID,
			&g.CreatedAt,
			&g.UpdateAt,
			&g.StudentID,
			&g.FirstName,
			&g.LastName,
			&g.DisciplineID,
			&g.DisciplineName,
			&g.Grade,
			&g.Comment,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, g)
	}
	return items, nil
}

// Средний балл по студенту/предмету с фильтрацией по датам
func (r *gradeJournalRepository) GetAverageGrade(
	ctx context.Context,
	studentID, disciplineID *int64,
	fromDate, toDate *time.Time,
) (float64, error) {
	query := `SELECT AVG(grade) FROM grade_journal WHERE 1=1`
	var args []interface{}
	if studentID != nil {
		query += " AND student_id = ?"
		args = append(args, *studentID)
	}
	if disciplineID != nil {
		query += " AND discipline_id = ?"
		args = append(args, *disciplineID)
	}
	if fromDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *fromDate)
	}
	if toDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *toDate)
	}
	var avg sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg.Valid {
		return avg.Float64, nil
	}
	return 0, nil
}
