package user

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type UserPosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *UserPosgresql {
	return &UserPosgresql{
		db: db,
	}
}
