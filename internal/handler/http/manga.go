package http

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

type createMangaRequest struct {
	Title          string                `form:"title"`
	SecondaryTitle *string               `form:"secondary_title"`
	Description    string                `form:"description"`
	Type           uint8                 `form:"type"`
	Status         uint8                 `form:"status"`
	AgeRestrict    uint8                 `form:"age_restrict"`
	ReleaseYear    uint16                `form:"release_year"`
	PreviewFile    *multipart.FileHeader `form:"preview_file"`
	Tags           string                `form:"tags"`
}

func (h *Handler) createManga(ctx echo.Context) error {
	var reqInput createMangaRequest

	err := ctx.Bind(&reqInput)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input"))
	}

	reqInput.PreviewFile, err = ctx.FormFile("preview_file")
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	h.logger.Debug("create manga request", map[string]any{"input": reqInput})

	if reqInput.PreviewFile == nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("preview file is required"))
	}

	if err := models.ValidateExtension(reqInput.PreviewFile.Filename, ".png", ".jpg", ".jpeg"); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	if reqInput.PreviewFile.Size > models.MaxMangaPreviewFileSize {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("preview file is too big"))
	}

	file, err := reqInput.PreviewFile.Open()
	if err != nil {
		return h.newAppErrorResponse(ctx,
			apperror.NewAppError(
				fmt.Errorf("failed to open preview file: %w", err),
				"Manga",
				"createManga",
				map[string]any{"file": reqInput.PreviewFile.Filename},
			))
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return h.newAppErrorResponse(ctx,
			apperror.NewAppError(
				fmt.Errorf("failed to read preview file: %w", err),
				"Manga",
				"createManga",
				map[string]any{"file": reqInput.PreviewFile.Filename},
			))
	}

	tags := strings.Split(
		strings.ReplaceAll(reqInput.Tags, " ", ""),
		",",
	)

	input, err := models.NewCreateMangaInput(
		reqInput.Title,
		reqInput.SecondaryTitle,
		reqInput.Description,
		models.NovelType(reqInput.Type),
		models.MangaStatus(reqInput.Status),
		models.AgeRestrict(reqInput.AgeRestrict),
		reqInput.ReleaseYear,
		models.UploadFile{
			Filename: reqInput.PreviewFile.Filename,
			Data:     fileBytes,
		},
		tags,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = h.services.Manga.Create(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaTitleExists):
			return h.newErrorResponse(ctx, http.StatusConflict, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

type getAllMangaResponse struct {
	Manga      []models.MangaOutput  `json:"manga"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getAllManga(ctx echo.Context) error {
	var filters models.MangaGetAllFilters

	err := ctx.Bind(&filters)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid filters"))
	}

	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	manga, count, err := h.services.Manga.GetAll(ctx.Request().Context(), models.DBPagination{
		Limit:  reqPagination.Limit(),
		Offset: reqPagination.Offset(),
	}, filters)
	if err != nil {
		switch {
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getAllMangaResponse{
		Manga:      manga,
		Pagination: reqPagination.GetFull(count),
	})

	return nil
}

type getMangaResponse struct {
	Manga models.MangaOutput `json:"manga"`
}

func (h *Handler) getManga(ctx echo.Context) error {
	var filters models.MangaFilters

	err := ctx.Bind(&filters)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid filters"))
	}

	manga, err := h.services.Manga.GetOne(ctx.Request().Context(), filters)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getMangaResponse{Manga: manga})

	return nil
}

type updateMangaRequest struct {
	Title          *string               `form:"title"`
	SecondaryTitle *string               `form:"secondary_title"`
	Description    *string               `form:"description"`
	Type           *models.NovelType     `form:"type"`
	Status         *models.MangaStatus   `form:"status"`
	AgeRestrict    *models.AgeRestrict   `form:"age_restrict"`
	ReleaseYear    *uint16               `form:"release_year"`
	PreviewFile    *multipart.FileHeader `form:"preview_file"`
	Tags           *string               `form:"tags"`
}

func (h *Handler) updateManga(ctx echo.Context) error {
	var reqInput updateMangaRequest

	err := ctx.Bind(&reqInput)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid manga ID"))
	}

	reqInput.PreviewFile, err = ctx.FormFile("preview_file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	h.logger.Debug("update manga request", map[string]any{"input": reqInput})

	var uploadFile *models.UploadFile
	if reqInput.PreviewFile != nil {
		if err := models.ValidateExtension(reqInput.PreviewFile.Filename, ".png", ".jpg", ".jpeg"); err != nil {
			return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
		}

		if reqInput.PreviewFile.Size > models.MaxMangaPreviewFileSize {
			return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("preview file is too big"))
		}
		file, err := reqInput.PreviewFile.Open()
		if err != nil {
			return h.newAppErrorResponse(ctx,
				apperror.NewAppError(
					fmt.Errorf("failed to open preview file: %w", err),
					"Manga",
					"updateManga",
					map[string]any{"file": reqInput.PreviewFile.Filename},
				))
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return h.newAppErrorResponse(ctx,
				apperror.NewAppError(
					fmt.Errorf("failed to read preview file: %w", err),
					"Manga",
					"updateManga",
					map[string]any{"file": reqInput.PreviewFile.Filename},
				))
		}

		uploadFile = &models.UploadFile{
			Filename: reqInput.PreviewFile.Filename,
			Data:     fileBytes,
		}
	}

	var tags []string

	if reqInput.Tags != nil && *reqInput.Tags != "" {
		tags = strings.Split(
			strings.ReplaceAll(*reqInput.Tags, " ", ""),
			",",
		)
	}

	input, err := models.NewUpdateMangaInput(
		id,
		reqInput.Title,
		reqInput.SecondaryTitle,
		reqInput.Description,
		reqInput.Type,
		reqInput.Status,
		reqInput.AgeRestrict,
		reqInput.ReleaseYear,
		uploadFile,
		&tags,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = h.services.Manga.Update(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaTitleExists):
			return h.newErrorResponse(ctx, http.StatusConflict, err.Error())
		case errors.Is(err, models.ErrTagNotFound):
			return h.newErrorResponse(ctx, http.StatusBadRequest, models.ErrTagNotFound.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusOK)

	return nil
}

func (h *Handler) deleteManga(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid manga ID"))
	}

	err = h.services.Manga.Delete(ctx.Request().Context(), id)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.NoContent(http.StatusOK)

	return nil
}
