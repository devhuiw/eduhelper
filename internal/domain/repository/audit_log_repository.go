package repository

import (
	"context"
	"database/sql"
	"service/internal/domain/models"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) AddAuditLog(ctx context.Context, entry *models.AuditLog) error {
	query := `INSERT INTO audit_log (user_id, table_name, row_id, action_type, old_data, new_data, comment)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		entry.UserID, entry.TableName, entry.RowID, entry.ActionType, entry.OldData, entry.NewData, entry.Comment)
	return err
}

func (r *AuditLogRepository) ListAuditLogs(ctx context.Context, limit, offset int) ([]*models.AuditLog, error) {
	query := `SELECT audit_id, created_at, user_id, table_name, row_id, action_type, old_data, new_data, comment
		FROM audit_log ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.AuditLog
	for rows.Next() {
		var a models.AuditLog
		err := rows.Scan(
			&a.AuditID, &a.CreatedAt, &a.UserID, &a.TableName, &a.RowID,
			&a.ActionType, &a.OldData, &a.NewData, &a.Comment,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}
