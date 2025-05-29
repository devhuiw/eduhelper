package models

import "time"

type Role struct {
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
	RoleName  string    `json:"role_name"`
}
