package service

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/service/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/service/manga"
	"github.com/lzaxel/zero-manga-backend/internal/service/user"
)

type User interface {
	Create(ctx context.Context, user models.CreateUserInput) error
}

type Manga interface{}

type Chapter interface{}

type Services struct {
	User
	Manga
	Chapter
}

func New(ctx context.Context, repository *repository.Repository) *Services {
	return &Services{
		User:    user.New(ctx, repository.User),
		Manga:   manga.New(ctx, repository.Manga),
		Chapter: chapter.New(ctx, repository.Chapter),
	}
}
