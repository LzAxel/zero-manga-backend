package manga

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type MangaPosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *MangaPosgresql {
	return &MangaPosgresql{
		db: db,
	}
}

func (m *MangaPosgresql) Create(ctx context.Context, manga models.Manga) error {
	query, args, _ := squirrel.
		Insert(postgresql.MangaTable).
		Columns(
			"id",
			"title",
			"secondary_title",
			"description",
			"slug",
			"type",
			"status",
			"age_restrict",
			"release_year",
			"preview_url",
		).
		Values(
			manga.ID,
			manga.Title,
			manga.SecondaryTitle,
			manga.Description,
			manga.Slug,
			manga.Type,
			manga.Status,
			manga.AgeRestrict,
			manga.ReleaseYear,
			manga.PreviewURL,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrMangaTitleExists
			}
		}
		return apperror.NewDBError(
			err,
			"Manga",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (m *MangaPosgresql) GetOne(ctx context.Context, filters models.MangaFilters) (models.Manga, error) {
	query := squirrel.
		Select("*").
		From(postgresql.MangaTable)

	if filters.ID != nil {
		query = query.Where(squirrel.Eq{"id": filters.ID})
	}
	if filters.Title != nil {
		query = query.Where(
			squirrel.Or{
				squirrel.Eq{"title": filters.Title},
				squirrel.Eq{"secondary_title": filters.Title},
			})
	}
	if filters.Slug != nil {
		query = query.Where(squirrel.Eq{"slug": filters.Slug})
	}

	queryString, args, _ := query.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var manga models.Manga
	if err := m.db.GetContext(ctx, &manga, queryString, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return manga, apperror.ErrNotFound
		}
		return manga, apperror.NewDBError(
			err,
			"Manga",
			"GetOne",
			queryString,
			args,
		)
	}

	return manga, nil
}

func (m *MangaPosgresql) GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.Manga, uint64, error) {
	// getting manga
	query := squirrel.
		Select("*").
		From(postgresql.MangaTable)

	queryString, args, _ := queryGetAllFilters(query, filters).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var count uint64
	var users = make([]models.Manga, 0)
	if err := m.db.SelectContext(ctx, &users, queryString, args...); err != nil {
		return users, count, apperror.NewDBError(
			err,
			"Manga",
			"GetAll",
			queryString,
			args,
		)
	}
	// counting users
	query = squirrel.
		Select("COUNT(*)").
		From(postgresql.MangaTable)

	queryString, args, _ = queryGetAllFilters(query, filters).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := m.db.GetContext(ctx, &count, queryString, args...); err != nil {
		return users, count, apperror.NewDBError(
			err,
			"Manga",
			"GetAll",
			queryString,
			args,
		)
	}

	return users, count, nil
}

func queryGetAllFilters(query squirrel.SelectBuilder, filters models.MangaGetAllFilters) squirrel.SelectBuilder {
	if len(filters.Type) > 0 {
		args := make([]interface{}, len(filters.Type))
		for i, v := range filters.Type {
			args[i] = v
		}
		query = query.Where(postgresql.FormatINClause("type", len(filters.Type)), args...)
	}

	if len(filters.Status) > 0 {
		args := make([]interface{}, len(filters.Status))
		for i, v := range filters.Status {
			args[i] = v
		}
		query = query.Where(postgresql.FormatINClause("status", len(filters.Status)), args...)
	}

	if len(filters.AgeRestrict) > 0 {
		args := make([]interface{}, len(filters.AgeRestrict))
		for i, v := range filters.AgeRestrict {
			args[i] = v
		}
		query = query.Where(postgresql.FormatINClause("age_restrict", len(filters.AgeRestrict)), args...)
	}

	if filters.ReleaseYear != nil {
		query = query.Where(squirrel.Eq{"release_year": filters.ReleaseYear})
	}

	return query
}

func (m *MangaPosgresql) Update(ctx context.Context, manga models.UpdateMangaRecord) error {
	query := squirrel.Update(postgresql.MangaTable)

	if manga.Title != nil {
		query = query.
			Set("title", *manga.Title).
			Set("slug", *manga.Slug)
	}
	if manga.SecondaryTitle != nil {
		query = query.Set("secondary_title", *manga.SecondaryTitle)
	}
	if manga.Description != nil {
		query = query.Set("description", *manga.Description)
	}
	if manga.Type != nil {
		query = query.Set("type", *manga.Type)
	}
	if manga.Status != nil {
		query = query.Set("status", *manga.Status)
	}
	if manga.AgeRestrict != nil {
		query = query.Set("age_restrict", *manga.AgeRestrict)
	}
	if manga.ReleaseYear != nil {
		query = query.Set("release_year", *manga.ReleaseYear)
	}
	if manga.PreviewURL != nil {
		query = query.Set("preview_url", *manga.PreviewURL)
	}

	queryString, args, _ := query.
		Where(squirrel.Eq{"id": manga.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, queryString, args...); err != nil {
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrMangaTitleExists
			}
		}
		return apperror.NewDBError(
			err,
			"Manga",
			"Update",
			queryString,
			args,
		)
	}

	return nil
}

func (m *MangaPosgresql) Delete(ctx context.Context, id uuid.UUID) error {
	queryString, args, _ := squirrel.
		Delete(postgresql.MangaTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, queryString, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Manga",
			"Delete",
			queryString,
			args,
		)
	}

	return nil
}
