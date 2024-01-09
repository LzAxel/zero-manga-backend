package auth

import (
	"context"
	"errors"
	"fmt"

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

func (a *Authorization) RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error) {
	claims, err := a.jwt.ValidateToken(refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken) || errors.Is(err, jwt.ErrInvalidClaims):
			return jwt.TokenPair{}, jwt.ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenExpired):
			return jwt.TokenPair{}, jwt.ErrTokenExpired
		}
		return jwt.TokenPair{}, fmt.Errorf("Authorization.RefreshTokens: %w", err)
	}

	tokenPair, err := a.jwt.GeneratePair(claims.Subject)
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("Authorization.RefreshTokens: %w", err)
	}

	return tokenPair, err
}

func (a *Authorization) Login(ctx context.Context, input models.LoginUserInput) (jwt.TokenPair, error) {
	user, err := a.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.As(err, &apperror.DBError{}) {
			dbErr := err.(apperror.DBError)
			if errors.Is(dbErr.Err, apperror.ErrNotFound) {
				return jwt.TokenPair{}, models.ErrInvalidCredentials
			}
		}

		return jwt.TokenPair{}, err
	}

	if err := hash.Compare(user.PasswordHash, input.Password); err != nil {
		return jwt.TokenPair{}, models.ErrInvalidCredentials
	}

	tokenPair, err := a.jwt.GeneratePair(user.ID)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	return tokenPair, err
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
