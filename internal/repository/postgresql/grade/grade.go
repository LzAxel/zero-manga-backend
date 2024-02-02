package grade

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

type GradePosgresql struct {
	db postgresql.DB
}

func New(db postgresql.DB) *GradePosgresql {
	return &GradePosgresql{
		db: db,
	}
}

func (p *GradePosgresql) Create(ctx context.Context, grade models.CreateGrade) error {
	query, args, _ := squirrel.
		Insert(postgresql.GradeTable).
		Columns(
			"user_id",
			"manga_id",
			"grade",
			"created_at").
		Values(
			grade.UserID,
			grade.MangaID,
			grade.Grade,
			grade.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := postgresql.GetPgError(err)
		if pgErr != nil {
			switch {
			case pgErr.Code == pgerrcode.UniqueViolation:
				return models.ErrDuplicatedGrade
			}
		}

		return apperror.NewDBError(
			err,
			"Grade",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (p *GradePosgresql) Delete(ctx context.Context, id int64) error {
	query, args, _ := squirrel.
		Delete(postgresql.GradeTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Grade",
			"Delete",
			query,
			args,
		)
	}

	return nil
}

type grade struct {
	ID        int64     `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	MangaID   uuid.UUID `db:"manga_id"`
	Grade     uint8     `db:"grade"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *GradePosgresql) GetByID(ctx context.Context, id int64) (models.Grade, error) {
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.GradeTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var grade grade
	if err := p.db.GetContext(ctx, &grade, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Grade{}, apperror.ErrNotFound
		default:
			return models.Grade{}, apperror.NewDBError(
				err,
				"Grade",
				"GetByID",
				query,
				args,
			)
		}
	}

	return models.NewGrade(
		grade.ID,
		grade.UserID,
		grade.MangaID,
		models.GradeType(grade.Grade),
		grade.CreatedAt,
	), nil
}

type getInfo struct {
	AvgGrade float64 `db:"avg_grade"`
	Count    uint64  `db:"count"`
}

func (p *GradePosgresql) GetInfoByManga(ctx context.Context, mangaID uuid.UUID) (float64, uint64, error) {
	query, args, _ := squirrel.
		Select("COALESCE(AVG(grade),0) as avg_grade, COUNT(*) as count").
		From(postgresql.GradeTable).
		Where(squirrel.Eq{"manga_id": mangaID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var info getInfo
	if err := p.db.GetContext(ctx, &info, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, 0, apperror.ErrNotFound
		default:
			return 0, 0, apperror.NewDBError(
				err,
				"Grade",
				"GetInfoByManga",
				query,
				args,
			)
		}
	}

	return info.AvgGrade, info.Count, nil
}

func (p *GradePosgresql) GetAllByUserID(ctx context.Context, pagination models.DBPagination, userID uuid.UUID) ([]models.Grade, uint64, error) {
	// getting
	query, args, _ := squirrel.
		Select("*").
		From(postgresql.GradeTable).
		Where(squirrel.Eq{"user_id": userID}).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var grades = make([]grade, 0)
	var count uint64

	if err := p.db.SelectContext(ctx, &grades, query, args...); err != nil {
		return []models.Grade{}, count, apperror.NewDBError(
			err,
			"Grade",
			"GetAllByUserID",
			query,
			args,
		)
	}

	gradesWithManga := make([]models.Grade, len(grades))
	for i, grade := range grades {
		gradesWithManga[i] = models.NewGrade(
			grade.ID,
			grade.UserID,
			grade.MangaID,
			models.GradeType(grade.Grade),
			grade.CreatedAt,
		)
	}

	// counting
	query, args, _ = squirrel.
		Select("COUNT(*)").
		From(postgresql.GradeTable).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := p.db.GetContext(ctx, &count, query, args...); err != nil {
		return []models.Grade{}, count, apperror.NewDBError(
			err,
			"Grade",
			"GetAllByUserID",
			query,
			args,
		)
	}

	return gradesWithManga, count, nil
}
