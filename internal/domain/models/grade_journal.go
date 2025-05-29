package models

import "time"

type GradeJournal struct {
	GradeJournalID int64     `json:"grade_journal_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"update_at"`
	StudentID      int64     `json:"student_id"`
	Grade          int16     `json:"grade"`
	Comment        *string   `json:"comment,omitempty"`
	DisciplineID   int64     `json:"discipline_id"`
}
