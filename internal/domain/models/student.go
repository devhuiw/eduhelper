package models

import "time"

type Student struct {
	UserID         int64     `json:"user_id"`
	Phone          string    `json:"phone"`
	Birthday       time.Time `json:"birthday"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"update_at"`
	StudentGroupID int64     `json:"student_group_id"`
}
