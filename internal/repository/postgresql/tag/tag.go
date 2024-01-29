package tag

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type TagPostresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *TagPostresql {
	return &TagPostresql{db: db}
}

func (t *TagPostresql) Create(ctx context.Context, tag models.Tag) error {
	query, args, _ := squirrel.
		Insert(postgresql.TagsTable).
		Columns(
			"id",
			"name",
			"slug",
			"is_nsfw",
			"created_at",
		).
		Values(
			tag.ID,
			tag.Name,
			tag.Slug,
			tag.IsNSFW,
			tag.CreatedAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := t.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrTagExists
			}
		}

		return apperror.NewDBError(
			err,
			"Tag",
			"Create",
			query,
			args,
		)
	}
	return nil
}

func (t *TagPostresql) GetAll(ctx context.Context) ([]models.Tag, error) {
	query, _, _ := squirrel.
		Select("*").
		From(postgresql.TagsTable).
		ToSql()
	var tags = make([]models.Tag, 0)
	if err := t.db.SelectContext(ctx, &tags, query); err != nil {
		return tags, apperror.NewDBError(
			err,
			"Tags",
			"GetAll",
			query,
			nil,
		)
	}
	return tags, nil
}

func (t *TagPostresql) Update(ctx context.Context, tag models.UpdateTagRecord) error {
	query := squirrel.Update(postgresql.TagsTable)

	if tag.Name != nil {
		query = query.
			Set("name", *tag.Name).
			Set("slug", *tag.Slug)
	}
	if tag.IsNSFW != nil {
		query = query.Set("is_nsfw", *tag.IsNSFW)
	}

	queryString, args, _ := query.
		Where(squirrel.Eq{"id": tag.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if _, err := t.db.ExecContext(ctx, queryString, args...); err != nil {
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrTagExists
			}
		}
		return apperror.NewDBError(
			err,
			"Tag",
			"Update",
			queryString,
			args,
		)
	}
	return nil
}

func (t *TagPostresql) Delete(ctx context.Context, id uuid.UUID) error {
	queryString, args, _ := squirrel.
		Delete(postgresql.TagsTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := t.db.ExecContext(ctx, queryString, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Tag",
			"Delete",
			queryString,
			args,
		)
	}
	return nil
}
