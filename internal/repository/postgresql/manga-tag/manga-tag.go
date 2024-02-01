package mangatag

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type MangaTagPostresql struct {
	db postgresql.DB
}

func New(db postgresql.DB) *MangaTagPostresql {
	return &MangaTagPostresql{db: db}
}

func (m *MangaTagPostresql) GetAllByMangaID(ctx context.Context, mangaID uuid.UUID) ([]models.MangaTagRelation, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.MangaTagsTable).
		Where(squirrel.Eq{"manga_id": mangaID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var tags = make([]models.MangaTagRelation, 0)
	if err := m.db.SelectContext(ctx, &tags, query, args...); err != nil {
		return tags, apperror.NewDBError(
			err,
			"MangaTag",
			"GetAllByMangaID",
			query,
			args,
		)
	}

	return tags, nil

}
func (m *MangaTagPostresql) Create(ctx context.Context, tag models.MangaTagRelation) error {
	query, args, _ := squirrel.
		Insert(postgresql.MangaTagsTable).
		Columns("manga_id", "tag_id").
		Values(tag.MangaID, tag.TagID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := postgresql.GetPgError(err)

		switch {
		case pgErr != nil && pgErr.Code == pgerrcode.ForeignKeyViolation:
			return models.ErrMangaOrTagNotFound
		default:
			return apperror.NewDBError(
				err,
				"MangaTag",
				"Create",
				query,
				args,
			)
		}
	}

	return nil
}
func (m *MangaTagPostresql) Delete(ctx context.Context, tag models.MangaTagRelation) error {
	query, args, _ := squirrel.
		Delete(postgresql.MangaTagsTable).
		Where(squirrel.Eq{"manga_id": tag.MangaID, "tag_id": tag.TagID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			err,
			"MangaTag",
			"Delete",
			query,
			args,
		)
	}

	return nil
}
