package models

import "time"

type User struct {
	UserID     int64     `json:"user_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdateAt   time.Time `json:"update_at,omitempty"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName *string   `json:"middle_name,omitempty"`
	Email      string    `json:"email"`
	Password   []byte    `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name,omitempty"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
}
