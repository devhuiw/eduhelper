package models

import "time"

type Teacher struct {
	UserID            int64     `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"updated_at"`
	Phone             string    `json:"phone"`
	WorkingExperience *string   `json:"working_experience,omitempty"`
	Education         *string   `json:"education,omitempty"`
}

type TeacherResponse struct {
	UserID            int64   `json:"user_id"`
	Phone             string  `json:"phone"`
	WorkingExperience *string `json:"working_experience,omitempty"`
	Education         *string `json:"education,omitempty"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	MiddleName        *string `json:"middle_name,omitempty"`
}

type TeacherPublic struct {
	UserID            int64   `json:"user_id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	MiddleName        *string `json:"middle_name,omitempty"`
	Education         string  `json:"education"`
	WorkingExperience *string `json:"working_experience,omitempty"`
}
