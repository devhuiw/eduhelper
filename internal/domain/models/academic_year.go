package models

import "time"

type AcademicYear struct {
	AcademicYearID int64     `json:"academic_year_id"`
	Name           string    `json:"name_academic_year"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
	StartWith      time.Time `json:"start_with"`
	EndsWith       time.Time `json:"ends_with"`
}
