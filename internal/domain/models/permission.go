package models

import "time"

type Permission struct {
	PermissionID   int64     `json:"permission_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
	PermissionName string    `json:"permission_name"`
}
