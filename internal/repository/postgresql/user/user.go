package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
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
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrUsernameEmailExists
			}
		}
		return apperror.NewDBError(
			postgresql.HandleDBError(err),
			"User",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (p *UserPosgresql) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.UsersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var user models.User
	if err := p.db.GetContext(ctx, &user, query, args...); err != nil {
		return user, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"User",
			"GetByID",
			query,
			args,
		)
	}

	return user, nil
}
func (p *UserPosgresql) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.UsersTable).
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var user models.User
	if err := p.db.GetContext(ctx, &user, query, args...); err != nil {
		return user, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"User",
			"GetByUsername",
			query,
			args,
		)
	}

	return user, nil
}
func (p *UserPosgresql) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.UsersTable).
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var user models.User
	if err := p.db.GetContext(ctx, &user, query, args...); err != nil {
		return user, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"User",
			"GetByEmail",
			query,
			args,
		)
	}

	return user, nil
}
