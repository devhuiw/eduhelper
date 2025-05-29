package models

import "time"

type StudentGroup struct {
	StudentGroupID   int64     `json:"student_group_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdateAt         time.Time `json:"update_at"`
	StudentGroupName string    `json:"student_group_name"`
	CuratorID        int64     `json:"curator_id"`
	AcademicYearID   int64     `json:"academic_year_id"`
}
