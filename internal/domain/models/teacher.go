package models

import "time"

type Teacher struct {
	UserID            int64     `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"update_at"`
	Phone             string    `json:"phone"`
	WorkingExperience *string   `json:"working_experience,omitempty"`
	Education         *string   `json:"education,omitempty"`
}
