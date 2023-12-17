package chapter

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository"
)

type Chapter struct {
	repo repository.Chapter
}

func New(ctx context.Context, repository repository.Chapter) *Chapter {
	return &Chapter{
		repo: repository,
	}
}
