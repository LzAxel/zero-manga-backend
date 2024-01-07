package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/jwt"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/service/auth"
	"github.com/lzaxel/zero-manga-backend/internal/service/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/service/manga"
	"github.com/lzaxel/zero-manga-backend/internal/service/uploader"
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

type Manga interface {
	Create(ctx context.Context, manga models.CreateMangaInput) error
	GetOne(ctx context.Context, filters models.MangaFilters) (models.MangaOutput, error)
	GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.MangaOutput, uint64, error)
	Update(ctx context.Context, manga models.UpdateMangaInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Chapter interface{}

type Services struct {
	User
	Manga
	Chapter
	Authorization
}

func New(
	ctx context.Context,
	repository *repository.Repository,
	jwt *jwt.JWT,
	fileStorage filestorage.FileStorage,
) *Services {
	uploader := uploader.NewUploader(fileStorage)
	return &Services{
		User:          user.New(ctx, repository.User),
		Manga:         manga.New(ctx, repository.Manga, uploader),
		Chapter:       chapter.New(ctx, repository.Chapter),
		Authorization: auth.New(ctx, jwt, repository.User),
	}
}
