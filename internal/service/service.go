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
	"github.com/lzaxel/zero-manga-backend/internal/service/grade"
	"github.com/lzaxel/zero-manga-backend/internal/service/manga"
	"github.com/lzaxel/zero-manga-backend/internal/service/uploader"
	"github.com/lzaxel/zero-manga-backend/internal/service/user"
)

type Authorization interface {
	Login(ctx context.Context, input models.LoginUserInput) (jwt.TokenPair, error)
	Register(ctx context.Context, input models.CreateUserInput) error
	RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error)
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

type Chapter interface {
	Create(ctx context.Context, chapter models.CreateChapterInput) error
	GetAllByManga(ctx context.Context, pagination models.DBPagination, mangaID uuid.UUID) ([]models.Chapter, uint64, error)
	Get(ctx context.Context, chapterID uuid.UUID) (models.ChapterOutput, error)
}

type Grade interface {
	Create(ctx context.Context, input models.CreateGradeInput) error
	Delete(ctx context.Context, user grade.UserInfo, gradeID int64) error
	GetAllByUserID(ctx context.Context, pagination models.DBPagination, userID uuid.UUID) ([]grade.GradeWithManga, uint64, error)
}

type Services struct {
	User
	Manga
	Chapter
	Authorization
	Grade
}

func New(
	ctx context.Context,
	repository *repository.Repository,
	jwt *jwt.JWT,
	fileStorage filestorage.FileStorage,
) *Services {
	uploader := uploader.NewUploader(fileStorage)
	manga := manga.New(repository.Manga, repository.Chapter, repository.Grade, uploader)
	return &Services{
		User:          user.New(repository.User),
		Manga:         manga,
		Chapter:       chapter.New(repository.Chapter, repository.Page, repository.Manga, uploader),
		Authorization: auth.New(jwt, repository.User),
		Grade:         grade.New(repository.Grade, manga),
	}
}
