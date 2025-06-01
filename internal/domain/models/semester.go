package models

import "time"

type Semester struct {
	SemesterID     int64     `json:"semester_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
	StartWith      time.Time `json:"start_with"`
	EndsWith       time.Time `json:"ends_with"`
	AcademicYearID int64     `json:"academic_year_id"`
}
