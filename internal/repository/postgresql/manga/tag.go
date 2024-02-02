package manga

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

func (m *MangaPostgresql) AddTagsToManga(ctx context.Context, mangaID uuid.UUID, tagIDs ...uuid.UUID) error {
	statement := squirrel.Insert(postgresql.MangaTagsTable).
		Columns("manga_id", "tag_id")

	for _, tagID := range tagIDs {
		statement = statement.Values(mangaID, tagID)
	}

	query, args, _ := statement.
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

func (m *MangaPostgresql) RemoveTagsFromManga(ctx context.Context, mangaID uuid.UUID, tagIDs ...uuid.UUID) error {
	var condition = make(squirrel.Or, 0)

	for _, tagID := range tagIDs {
		condition = append(condition, squirrel.Eq{"manga_id": mangaID, "tag_id": tagID})
	}

	query, args, _ := squirrel.
		Delete(postgresql.MangaTagsTable).
		Where(condition).
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
