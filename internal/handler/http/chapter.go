package http

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

type createChapterRequest struct {
	MangaID      string                `form:"manga_id"`
	Title        *string               `form:"title"`
	Number       uint                  `form:"number"`
	Volume       uint                  `form:"volume"`
	PagesZipFile *multipart.FileHeader `form:"pages_zip_file"`
}

func (h *Handler) createChapter(ctx echo.Context) error {
	var reqInput createChapterRequest

	err := ctx.Bind(&reqInput)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	reqInput.PagesZipFile, err = ctx.FormFile("pages_zip_file")
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("pages zip file is required"))
	}

	h.logger.Debug("create chapter request", map[string]any{"input": reqInput})

	if reqInput.PagesZipFile == nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("pages zip file is required"))
	}

	if err := models.ValidateExtension(reqInput.PagesZipFile.Filename, ".zip"); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	if reqInput.PagesZipFile.Size > models.MaxMangaPageZipFileSize {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("pages zip file is too big"))
	}

	file, err := reqInput.PagesZipFile.Open()
	defer file.Close()
	if err != nil {
		return h.newAppErrorResponse(ctx,
			apperror.NewAppError(
				fmt.Errorf("failed to open pages zip file: %w", err),
				"Chapter",
				"createChapter",
				map[string]any{"file": reqInput.PagesZipFile.Filename},
			))
	}
	zipReader := io.LimitReader(file, models.MaxMangaPageZipReaderSize)

	mangaID, err := uuid.Parse(reqInput.MangaID)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid manga ID"))
	}
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	input := models.CreateChapterInput{
		MangaID:    mangaID,
		UploaderID: user.ID,
		Title:      reqInput.Title,
		Number:     reqInput.Number,
		Volume:     reqInput.Volume,
		PageArchive: models.UploadReader{
			Reader:   zipReader,
			Filename: reqInput.PagesZipFile.Filename,
		},
	}

	err = h.services.Chapter.Create(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, models.ErrMangaNotFound.Error())
		case errors.Is(err, models.ErrNoValidImages):
			return h.newValidationErrorResponse(ctx, http.StatusBadRequest, models.ErrNoValidImages)
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

type getAllChaptersResponse struct {
	Chapters   []models.Chapter      `json:"chapters"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getChapterByManga(ctx echo.Context) error {
	mangaID, err := uuid.Parse(ctx.Param("manga_id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid manga ID"))
	}

	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	chapters, count, err := h.services.Chapter.GetAllByManga(ctx.Request().Context(), models.DBPagination{
		Limit:  reqPagination.Limit(),
		Offset: reqPagination.Offset(),
	}, mangaID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, models.ErrMangaNotFound.Error())
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, getAllChaptersResponse{
		Chapters:   chapters,
		Pagination: reqPagination.GetFull(count),
	})

	return nil
}

func (h *Handler) getChapter(ctx echo.Context) error {
	mangaID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid manga ID"))
	}

	chapter, err := h.services.Chapter.Get(ctx.Request().Context(), mangaID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrChapterNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, models.ErrChapterNotFound.Error())
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, chapter)

	return nil
}
