package tag

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func New(db postgresql.DB) *TagPostresql {
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
				return models.ErrTagDuplicated
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
		PlaceholderFormat(squirrel.Dollar).
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

func (t *TagPostresql) GetByID(ctx context.Context, id uuid.UUID) (models.Tag, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.TagsTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var tag = models.Tag{}
	if err := t.db.GetContext(ctx, &tag, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return tag, apperror.ErrNotFound
		default:
			return tag, apperror.NewDBError(
				err,
				"Tag",
				"GetByID",
				query,
				args,
			)

		}
	}

	return tag, nil
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
				return models.ErrTagDuplicated
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

func (m *TagPostresql) GetAllByMangaID(ctx context.Context, mangaID uuid.UUID) ([]models.Tag, error) {
	query, args, _ := squirrel.
		Select("tag.*").
		From(fmt.Sprintf("%s manga_tag", postgresql.MangaTagsTable)).
		LeftJoin(fmt.Sprintf("%s tag ON manga_tag.tag_id = tag.id", postgresql.TagsTable)).
		Where(squirrel.Eq{"manga_id": mangaID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var tags = make([]models.Tag, 0)
	if err := m.db.SelectContext(ctx, &tags, query, args...); err != nil {
		return tags, apperror.NewDBError(
			err,
			"Tag",
			"GetAllByMangaID",
			query,
			args,
		)
	}

	return tags, nil

}
