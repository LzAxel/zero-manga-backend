package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
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
	return u.repo.GetByID(ctx, id)
}
func (u *User) GetByUsername(ctx context.Context, username string) (models.User, error) {
	return u.repo.GetByUsername(ctx, username)
}
func (u *User) GetByEmail(ctx context.Context, email string) (models.User, error) {
	return u.repo.GetByEmail(ctx, email)
}

func (u *User) GetAll(ctx context.Context, pagination models.Pagination, filters models.UserFilters) ([]models.User, models.FullPagination, error) {
	users, total, err := u.repo.GetAll(ctx, postgresql.Pagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	}, filters)

	return users, pagination.GetFull(total), err
}
