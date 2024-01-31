package grade

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
)

type MangaService interface {
	GetOne(ctx context.Context, filters models.MangaFilters) (models.MangaOutput, error)
}

type Grade struct {
	repo  repository.Grade
	manga MangaService
}

func New(repo repository.Grade, manga MangaService) *Grade {
	return &Grade{
		repo:  repo,
		manga: manga,
	}
}

func (g *Grade) Create(ctx context.Context, input models.CreateGradeInput) error {
	grade, err := models.NewCreateGrade(
		input.UserID,
		input.MangaID,
		input.Grade,
		clock.Now(),
	)
	if err != nil {
		return err
	}

	err = g.repo.Create(ctx, grade)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedGrade) {
			return models.ErrDuplicatedGrade
		}
		return err
	}

	return nil
}

type UserInfo struct {
	ID   uuid.UUID
	Type models.UserType
}

func (g *Grade) Delete(ctx context.Context, user UserInfo, gradeID int64) error {
	if user.Type != models.UserTypeAdmin {
		grade, err := g.repo.GetByID(ctx, gradeID)
		if err != nil {
			return handleNotFoundError(err)
		}
		if grade.UserID != user.ID {
			return models.ErrNotCreatorOfGrade
		}
	}
	err := g.repo.Delete(ctx, gradeID)
	if err != nil {
		return err
	}

	return nil
}

type GradeWithManga struct {
	ID        int64
	Manga     models.MangaOutput
	GradeType models.GradeType
	CreatedAt time.Time
}

func NewGradeWithManga(id int64, manga models.MangaOutput, gradeType models.GradeType, createdAt time.Time) GradeWithManga {
	return GradeWithManga{
		ID:        id,
		Manga:     manga,
		GradeType: gradeType,
		CreatedAt: createdAt,
	}
}

func (g *Grade) GetAllByUserID(ctx context.Context, pagination models.DBPagination, userID uuid.UUID) ([]GradeWithManga, uint64, error) {
	grades, count, err := g.repo.GetAllByUserID(ctx, pagination, userID)
	if err != nil {
		return nil, 0, err
	}

	var gradesWithManga = make([]GradeWithManga, len(grades))
	for i, grade := range grades {
		manga, err := g.manga.GetOne(ctx, models.MangaFilters{
			ID: &grade.MangaID,
		})
		if err != nil {
			return nil, 0, apperror.NewAppError(
				fmt.Errorf("failed to get manga for grade: %w", err),
				"Grade",
				"GetAllByUserID",
				map[string]interface{}{
					"iteration": i,
					"mangaID":   grade.MangaID,
					"userID":    userID,
				},
			)
		}
		gradesWithManga[i] = NewGradeWithManga(grade.ID, manga, grade.Grade, grade.CreatedAt)
	}

	return gradesWithManga, count, nil
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrGradeNotFound
	}
	return err
}
