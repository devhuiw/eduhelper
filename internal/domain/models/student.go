package models

import "time"

type Student struct {
	UserID         int64     `json:"user_id"`
	Phone          string    `json:"phone"`
	Birthday       time.Time `json:"birthday"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
	StudentGroupID int64     `json:"student_group_id"`
}

type StudentPublic struct {
	UserID         int64     `json:"user_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MiddleName     *string   `json:"middle_name,omitempty"`
	Birthday       time.Time `json:"birthday"`
	StudentGroupID int64     `json:"student_group_id"`
}
