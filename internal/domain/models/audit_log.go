package models

import "time"

type AuditLog struct {
	AuditID    int64     `json:"audit_id"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     *int64    `json:"user_id,omitempty"`
	TableName  string    `json:"table_name"`
	RowID      int64     `json:"row_id"`
	ActionType string    `json:"action_type"`
	OldData    *string   `json:"old_data,omitempty"`
	NewData    *string   `json:"new_data,omitempty"`
	Comment    *string   `json:"comment,omitempty"`
}
