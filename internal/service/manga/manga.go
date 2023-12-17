package manga

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository"
)

type Manga struct {
	repo repository.Manga
}

func New(ctx context.Context, repository repository.Manga) *Manga {
	return &Manga{
		repo: repository,
	}
}
