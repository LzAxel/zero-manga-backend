package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/grade"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/manga"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/page"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/user"
)

type User interface {
	Create(ctx context.Context, user models.CreateUserRecord) error
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetAll(ctx context.Context, pagination models.DBPagination, filters models.UserFilters) ([]models.User, uint64, error)
}

type Manga interface {
	Create(ctx context.Context, manga models.Manga) error
	GetOne(ctx context.Context, filters models.MangaFilters) (models.Manga, error)
	GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.Manga, uint64, error)
	Update(ctx context.Context, manga models.UpdateMangaRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Chapter interface {
	Create(ctx context.Context, chapter models.Chapter) error
	GetAllByManga(ctx context.Context, pagination models.DBPagination, mangaID uuid.UUID) ([]models.Chapter, uint64, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Chapter, error)
	GetByNumber(ctx context.Context, filter models.ChapterFilter) (models.Chapter, error)
	CountByManga(ctx context.Context, mangaID uuid.UUID) (uint64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Page interface {
	GetAllByChapter(ctx context.Context, chapterID uuid.UUID) ([]models.Page, error)
	Create(ctx context.Context, page models.Page) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Grade interface {
	Create(ctx context.Context, grade models.CreateGrade) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (models.Grade, error)
	GetInfoByManga(ctx context.Context, mangaID uuid.UUID) (float64, uint64, error)
	GetAllByUserID(ctx context.Context, pagination models.DBPagination, userID uuid.UUID) ([]models.Grade, uint64, error)
}

type Repository struct {
	User
	Manga
	Chapter
	Page
	Grade
}

func New(psql postgresql.PostgresqlRepository, logger logger.Logger) *Repository {
	return &Repository{
		User:    user.New(psql.DB),
		Manga:   manga.New(psql.DB),
		Chapter: chapter.New(psql.DB),
		Page:    page.New(psql.DB),
		Grade:   grade.New(psql.DB),
	}
}
