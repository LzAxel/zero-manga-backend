package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/service/grade"
)

type createGradeRequest struct {
	MangaID string `json:"manga_id"`
	Grade   uint8  `json:"grade"`
}

func (h *Handler) createGrade(ctx echo.Context) error {
	var input createGradeRequest

	if err := ctx.Bind(&input); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input"))
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	grade, err := models.NewCreateGradeInput(
		user.ID,
		input.MangaID,
		input.Grade,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	if err := h.services.Grade.Create(ctx.Request().Context(), grade); err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicatedGrade):
			return h.newErrorResponse(ctx, http.StatusConflict, models.ErrDuplicatedGrade.Error())
		case errors.Is(err, models.ErrInvalidGradeType):
			return h.newErrorResponse(ctx, http.StatusBadRequest, models.ErrInvalidGradeType.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

type getGradeInfoByMangaResponse struct {
	AvgGrade float64 `json:"avg_grade"`
	Count    uint64  `json:"count"`
}

func (h *Handler) deleteGrade(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	gradeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid grade id"))
	}

	if err := h.services.Grade.Delete(ctx.Request().Context(),
		grade.UserInfo{
			ID:   user.ID,
			Type: user.Type,
		},
		gradeID,
	); err != nil {
		switch {
		case errors.Is(err, models.ErrNotCreatorOfGrade):
			return h.newErrorResponse(
				ctx, http.StatusForbidden, models.ErrNotCreatorOfGrade.Error())
		case errors.Is(err, models.ErrGradeNotFound):
			return h.newErrorResponse(
				ctx, http.StatusNotFound, models.ErrGradeNotFound.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusOK)

	return nil
}

type gradeResponse struct {
	ID        int64            `json:"grade_id"`
	GradeType models.GradeType `json:"grade_type"`
	CreatedAt time.Time        `json:"created_at"`
}

type gradeWithManga struct {
	Manga models.MangaOutput `json:"manga"`
	Grade gradeResponse      `json:"grade"`
}

type getGradesByUserIDResponse struct {
	Grades     []gradeWithManga      `json:"grades"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getGradesByUserID(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid user ID"))
	}

	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	grades, count, err := h.services.Grade.GetAllByUserID(ctx.Request().Context(),
		models.DBPagination{
			Limit:  reqPagination.Limit(),
			Offset: reqPagination.Offset(),
		},
		userID,
	)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	var gradesWithManga = make([]gradeWithManga, len(grades))
	for i, grade := range grades {
		gradesWithManga[i] = gradeWithManga{
			Manga: grade.Manga,
			Grade: gradeResponse{
				ID:        grade.ID,
				GradeType: grade.GradeType,
				CreatedAt: grade.CreatedAt,
			},
		}
	}

	ctx.JSON(http.StatusOK, getGradesByUserIDResponse{
		Grades:     gradesWithManga,
		Pagination: reqPagination.GetFull(count),
	})

	return nil
}
