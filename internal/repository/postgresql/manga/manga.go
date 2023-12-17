package manga

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type MangaPosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *MangaPosgresql {
	return &MangaPosgresql{
		db: db,
	}
}
