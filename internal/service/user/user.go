package user

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
	"github.com/lzaxel/zero-manga-backend/pkg/hash"
	"github.com/lzaxel/zero-manga-backend/pkg/uuid"
)

type User struct {
	repo repository.User
}

func New(ctx context.Context, repository repository.User) *User {
	return &User{
		repo: repository,
	}
}

func (u *User) Create(ctx context.Context, user models.CreateUserInput) error {
	passwordHash, err := hash.Hash(user.Password)
	if err != nil {
		return err
	}
	dto := models.CreateUserRecord{
		ID:           uuid.New(),
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		PasswordHash: passwordHash,
		Gender:       int(user.Gender),
		Bio:          user.Bio,
		Type:         models.UserTypeReader,
		OnlineAt:     clock.Now(),
		RegisteredAt: clock.Now(),
	}

	return u.repo.Create(ctx, dto)
}
