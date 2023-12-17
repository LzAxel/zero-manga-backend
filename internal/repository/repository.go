package repository

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/chapter"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/manga"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql/user"
)

type User interface{}

type Manga interface{}

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
