package repository

import (
	"context"
	"service/internal/domain/models"
)

type GradeJournalRepository interface {
	Create(ctx context.Context, gj *models.GradeJournal) error
	GetByID(ctx context.Context, id int64) (*models.GradeJournal, error)
	Update(ctx context.Context, gj *models.GradeJournal) error
	Delete(ctx context.Context, id int64) error
	ListByStudent(ctx context.Context, studentID int64, limit, offset int) ([]*models.GradeJournal, error)
}
