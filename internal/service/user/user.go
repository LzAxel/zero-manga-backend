package user

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/repository"
)

type User struct {
	repo repository.User
}

func New(ctx context.Context, repository repository.User) *User {
	return &User{
		repo: repository,
	}
}
