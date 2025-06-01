package models

import "time"

type Discipline struct {
	DisciplineID   int64     `json:"discipline_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
	DisciplineName string    `json:"discipline_name"`
	TeacherID      int64     `json:"teacher_id"`
	StudentGroupID int64     `json:"student_group_id"`
}

type DisciplinePublic struct {
	DisciplineID      int64     `json:"discipline_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"updated_at"`
	DisciplineName    string    `json:"discipline_name"`
	TeacherID         int64     `json:"teacher_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	MiddleName        *string   `json:"middle_name,omitempty"`
	StudentGroupID    int64     `json:"student_group_id"`
	StudentGroupName  string    `json:"student_group_name"`
	CuratorID         int64     `json:"curator_id"`
	CuratorFirstName  string    `json:"curator_first_name"`
	CuratorLastName   string    `json:"curator_last_name"`
	CuratorMiddleName *string   `json:"curator_middle_name,omitempty"`
	AcademicYearID    int64     `json:"academic_year_id"`
}
