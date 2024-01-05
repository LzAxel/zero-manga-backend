package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/jwt"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/service/auth"
	"github.com/lzaxel/zero-manga-backend/internal/service/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/service/manga"
	"github.com/lzaxel/zero-manga-backend/internal/service/user"
)

type Authorization interface {
	Login(ctx context.Context, input models.LoginUserInput) (string, error)
	Register(ctx context.Context, input models.CreateUserInput) error
}

type User interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetAll(ctx context.Context, pagination models.Pagination, filters models.UserFilters) ([]models.User, models.FullPagination, error)
}

type Manga interface{}

type Chapter interface{}

type Services struct {
	User
	Manga
	Chapter
	Authorization
}

func New(ctx context.Context, repository *repository.Repository, jwt *jwt.JWT) *Services {
	return &Services{
		User:          user.New(ctx, repository.User),
		Manga:         manga.New(ctx, repository.Manga),
		Chapter:       chapter.New(ctx, repository.Chapter),
		Authorization: auth.New(ctx, jwt, repository.User),
	}
}
