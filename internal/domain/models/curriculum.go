package models

import "time"

type Curriculum struct {
	CurriculumID       int64     `json:"curriculum_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdateAt           time.Time `json:"updated_at"`
	SubjectName        string    `json:"subject_name"`
	SubjectDescription *string   `json:"subject_description,omitempty"`
	SemesterID         *int64    `json:"semester_id,omitempty"`
	DisciplineID       int64     `json:"discipline_id"`
}
