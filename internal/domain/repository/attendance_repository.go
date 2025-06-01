package repository

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/domain/models"
	"time"
)

type attendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *attendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) CreateAttendance(ctx context.Context, a *models.Attendance) error {
	query := `
		INSERT INTO attendance (created_at, visit, comment, updated_at, student_id, discipline_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	a.CreatedAt = now
	a.UpdateAt = now
	res, err := r.db.ExecContext(ctx, query, a.CreatedAt, a.Visit, a.Comment, a.UpdateAt, a.StudentID, a.DisciplineID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		a.AttendanceID = id
	}
	return err
}

func (r *attendanceRepository) GetAttendanceByID(ctx context.Context, id int64) (*models.Attendance, error) {
	query := `
		SELECT attendance_id, created_at, visit, comment, updated_at, student_id, discipline_id
		FROM attendance
		WHERE attendance_id = ?
	`
	a := &models.Attendance{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.AttendanceID,
		&a.CreatedAt,
		&a.Visit,
		&a.Comment,
		&a.UpdateAt,
		&a.StudentID,
		&a.DisciplineID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return a, nil
}

func (r *attendanceRepository) UpdateAttendance(ctx context.Context, a *models.Attendance) error {
	query := `
		UPDATE attendance
		SET visit = ?, comment = ?, updated_at = ?, student_id = ?, discipline_id = ?
		WHERE attendance_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, a.Visit, a.Comment, time.Now(), a.StudentID, a.DisciplineID, a.AttendanceID)
	return err
}

func (r *attendanceRepository) DeleteAttendance(ctx context.Context, id int64) error {
	query := `DELETE FROM attendance WHERE attendance_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *attendanceRepository) ListAttendance(ctx context.Context, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT attendance_id, created_at, visit, comment, updated_at, student_id, discipline_id
		FROM attendance
		ORDER BY attendance_id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Attendance
	for rows.Next() {
		a := &models.Attendance{}
		err := rows.Scan(
			&a.AttendanceID,
			&a.CreatedAt,
			&a.Visit,
			&a.Comment,
			&a.UpdateAt,
			&a.StudentID,
			&a.DisciplineID,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

func (r *attendanceRepository) ListAttendanceWithFilters(
	ctx context.Context,
	studentID, disciplineID *int64,
	date *time.Time,
	limit, offset int,
) ([]*models.Attendance, error) {
	query := `SELECT attendance_id, created_at, visit, comment, updated_at, student_id, discipline_id FROM attendance WHERE 1=1`
	var args []interface{}

	if studentID != nil {
		query += " AND student_id = ?"
		args = append(args, *studentID)
	}
	if disciplineID != nil {
		query += " AND discipline_id = ?"
		args = append(args, *disciplineID)
	}
	if date != nil {
		query += " AND DATE(created_at) = ?"
		args = append(args, date.Format("2006-01-02"))
	}
	query += " ORDER BY attendance_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Attendance
	for rows.Next() {
		a := &models.Attendance{}
		err := rows.Scan(
			&a.AttendanceID,
			&a.CreatedAt,
			&a.Visit,
			&a.Comment,
			&a.UpdateAt,
			&a.StudentID,
			&a.DisciplineID,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}
