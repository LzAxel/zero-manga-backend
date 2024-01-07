package page

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type PagePosgresql struct {
	db postgresql.DB
}

func New(ctx context.Context, db postgresql.DB) *PagePosgresql {
	return &PagePosgresql{
		db: db,
	}
}

func (p *PagePosgresql) GetAllByChapter(ctx context.Context, chapterID uuid.UUID) ([]models.Page, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.PageTable).
		Where(squirrel.Eq{"chapter_id": chapterID}).
		OrderBy("number ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var pages []models.Page
	if err := p.db.SelectContext(ctx, &pages, query, args...); err != nil {
		return nil, apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Page",
			"GetAllByChapter",
			query,
			args,
		)
	}

	return pages, nil
}
func (p *PagePosgresql) Create(ctx context.Context, page models.Page) error {
	query, args, _ := squirrel.
		Insert(postgresql.PageTable).
		Columns(
			"id",
			"chapter_id",
			"url",
			"number",
			"height",
			"width",
			"created_at").
		Values(
			page.ID,
			page.ChapterID,
			page.URL,
			page.Number,
			page.Height,
			page.Width,
			page.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Page",
			"Create",
			query,
			args,
		)
	}

	return nil
}
func (p *PagePosgresql) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, _ := squirrel.
		Delete(postgresql.PageTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			postgresql.HandleDBError(err),
			"Page",
			"Delete",
			query,
			args,
		)
	}

	return nil
}
