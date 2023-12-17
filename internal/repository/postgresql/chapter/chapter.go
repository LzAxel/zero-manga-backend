package chapter

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type ChapterPosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *ChapterPosgresql {
	return &ChapterPosgresql{
		db: db,
	}
}
