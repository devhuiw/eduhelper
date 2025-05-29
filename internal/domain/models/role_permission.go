package models

import "time"

type RolePermission struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"update_at"`
	RoleID       int64     `json:"role_id"`
	PermissionID int64     `json:"permission_id"`
}
