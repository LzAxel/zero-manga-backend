package chapter

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type ChapterPosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *ChapterPosgresql {
	return &ChapterPosgresql{
		db: db,
	}
}

func (c *ChapterPosgresql) Create(ctx context.Context, chapter models.Chapter) error {
	query, args, _ := squirrel.
		Insert(postgresql.ChapterTable).
		Columns(
			"id",
			"manga_id",
			"title",
			"number",
			"volume",
			"page_count",
			"uploader_id",
			"uploaded_at",
		).
		Values(
			chapter.ID,
			chapter.MangaID,
			chapter.Title,
			chapter.Number,
			chapter.Volume,
			chapter.PageCount,
			chapter.UploaderID,
			chapter.UploadedAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := c.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (m *ChapterPosgresql) GetAllByManga(ctx context.Context, pagination models.DBPagination, mangaID uuid.UUID) ([]models.Chapter, uint64, error) {
	// getting chapter
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.ChapterTable).
		Where(squirrel.Eq{
			"manga_id": mangaID,
		}).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var count uint64
	var chapters = make([]models.Chapter, 0)
	if err := m.db.SelectContext(ctx, &chapters, query, args...); err != nil {
		return chapters, count, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"GetAllByManga",
			query,
			args,
		)
	}

	// counting chapters
	query, args, _ = squirrel.
		Select("COUNT(*)").
		Where(squirrel.Eq{
			"manga_id": mangaID,
		}).
		From(postgresql.ChapterTable).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := m.db.GetContext(ctx, &count, query, args...); err != nil {
		return chapters, count, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"GetAllByManga",
			query,
			args,
		)
	}

	return chapters, count, nil
}

func (m *ChapterPosgresql) CountByManga(ctx context.Context, mangaID uuid.UUID) (uint64, error) {
	var count uint64
	query, args, _ := squirrel.
		Select("COUNT(*)").
		Where(squirrel.Eq{
			"manga_id": mangaID,
		}).
		From(postgresql.ChapterTable).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := m.db.GetContext(ctx, &count, query, args...); err != nil {
		return count, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"GetAllByManga",
			query,
			args,
		)
	}

	return count, nil
}

func (c *ChapterPosgresql) GetByID(ctx context.Context, id uuid.UUID) (models.Chapter, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.ChapterTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var chapter models.Chapter
	if err := c.db.GetContext(ctx, &chapter, query, args...); err != nil {
		return models.Chapter{}, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"GetByID",
			query,
			args,
		)
	}

	return chapter, nil
}

func (c *ChapterPosgresql) GetByNumber(ctx context.Context, filter models.ChapterFilter) (models.Chapter, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.ChapterTable).
		Where(squirrel.Eq{
			"manga_id": filter.MangaID,
			"number":   filter.Number,
			"volume":   filter.Volume,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var chapter models.Chapter
	if err := c.db.GetContext(ctx, &chapter, query, args...); err != nil {
		return models.Chapter{}, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"GetByNumber",
			query,
			args,
		)
	}

	return chapter, nil
}

func (c *ChapterPosgresql) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, _ := squirrel.
		Delete(postgresql.ChapterTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := c.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Chapter",
			"Delete",
			query,
			args,
		)
	}

	return nil
}
