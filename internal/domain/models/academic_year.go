package models

import "time"

type AcademicYear struct {
	AcademicYearID int64     `json:"academic_year_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"update_at"`
	StartWith      time.Time `json:"start_with"`
	EndsWith       time.Time `json:"ends_with"`
}
