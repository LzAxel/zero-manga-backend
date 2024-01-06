package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/manga"
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

type Chapter interface{}

type Repository struct {
	User
	Manga
	Chapter
}

func New(ctx context.Context, psql postgresql.PostgresqlRepository, logger logger.Logger) *Repository {
	return &Repository{
		User:    user.New(ctx, psql.DB),
		Manga:   manga.New(ctx, psql.DB),
		Chapter: chapter.New(ctx, psql.DB),
	}
}
