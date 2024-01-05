package auth

import (
	"context"
	"errors"

	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/jwt"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
	"github.com/lzaxel/zero-manga-backend/pkg/hash"
	"github.com/lzaxel/zero-manga-backend/pkg/uuid"
)

type Authorization struct {
	jwt      *jwt.JWT
	userRepo repository.User
}

func New(ctx context.Context, jwt *jwt.JWT, userRepo repository.User) *Authorization {
	return &Authorization{
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (a *Authorization) Login(ctx context.Context, input models.LoginUserInput) (string, error) {
	user, err := a.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.As(err, &apperror.DBError{}) {
			dbErr := err.(apperror.DBError)
			if errors.Is(dbErr.Err, apperror.ErrNotFound) {
				return "", models.ErrInvalidCredentials
			}
		}

		return "", err
	}

	if err := hash.Compare(user.PasswordHash, input.Password); err != nil {
		return "", models.ErrInvalidCredentials
	}

	token, err := a.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
func (a *Authorization) Register(ctx context.Context, input models.CreateUserInput) error {
	passwordHash, err := hash.Hash(input.Password)
	if err != nil {
		return err
	}

	dto := models.CreateUserRecord{
		ID:           uuid.New(),
		Username:     input.Username,
		DisplayName:  input.DisplayName,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Gender:       int(input.Gender),
		Bio:          input.Bio,
		Type:         models.UserTypeReader,
		OnlineAt:     clock.Now(),
		RegisteredAt: clock.Now(),
	}

	return a.userRepo.Create(ctx, dto)
}
