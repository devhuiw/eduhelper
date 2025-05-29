package models

import "time"

type Discipline struct {
	DisciplineID   int64     `json:"discipline_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"update_at"`
	DisciplineName string    `json:"discipline_name"`
	TeacherID      int64     `json:"teacher_id"`
	StudentGroupID int64     `json:"student_group_id"`
}
