package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/lzaxel/zero-manga-backend/internal/models"
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

func (p *UserPosgresql) Create(ctx context.Context, user models.CreateUserRecord) error {
	query, args, _ := squirrel.
		Insert(postgresql.UsersTable).
		Columns(
			"id",
			"username",
			"display_name",
			"email",
			"password_hash",
			"gender",
			"bio",
			"type",
			"online_at",
			"registered_at",
		).
		Values(
			user.ID,
			user.Username,
			user.DisplayName,
			user.Email,
			user.PasswordHash,
			user.Gender,
			user.Bio,
			user.Type,
			user.OnlineAt,
			user.RegisteredAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
