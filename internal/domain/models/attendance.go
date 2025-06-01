package models

import "time"

type Attendance struct {
	AttendanceID int64     `json:"attendance_id"`
	CreatedAt    time.Time `json:"created_at"`
	Visit        bool      `json:"visit"`
	Comment      *string   `json:"comment,omitempty"`
	UpdateAt     time.Time `json:"updated_at"`
	StudentID    int64     `json:"student_id"`
	DisciplineID int64     `json:"discipline_id"`
}
