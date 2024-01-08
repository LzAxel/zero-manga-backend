package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
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

func (u *User) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, err := u.repo.GetByID(ctx, id)

	return user, handleNotFoundError(err)
}
func (u *User) GetByUsername(ctx context.Context, username string) (models.User, error) {
	user, err := u.repo.GetByUsername(ctx, username)

	return user, handleNotFoundError(err)
}
func (u *User) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)

	return user, handleNotFoundError(err)
}

func (u *User) GetAll(ctx context.Context, pagination models.Pagination, filters models.UserFilters) ([]models.User, models.FullPagination, error) {
	users, total, err := u.repo.GetAll(ctx, models.DBPagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	}, filters)

	return users, pagination.GetFull(total), err
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrUserNotFound
	}
	return err
}
