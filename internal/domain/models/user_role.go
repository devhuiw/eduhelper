package models

import "time"

type UserRole struct {
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
	RoleID    int64     `json:"role_id"`
	UserID    int64     `json:"user_id"`
}
