package http

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"net/http"
)

type createTagRequest struct {
	Name   string `json:"name"`
	IsNSFW bool   `json:"is_nsfw"`
}

func (h *Handler) createTag(ctx echo.Context) error {
	var reqInput createTagRequest

	err := ctx.Bind(&reqInput)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	input, err := models.NewCreateTagInput(
		reqInput.Name,
		reqInput.IsNSFW,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	err = h.services.Tag.Create(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTagDuplicated):
			return h.newErrorResponse(ctx, http.StatusConflict, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}
	ctx.NoContent(http.StatusCreated)
	return nil
}

type getAllTagsResponse struct {
	Tags []models.Tag `json:"tags"`
}

func (h *Handler) getAllTags(ctx echo.Context) error {
	tags, err := h.services.Tag.GetAll(ctx.Request().Context())
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, getAllTagsResponse{
		Tags: tags,
	})
	return nil
}

type updateTagRequest struct {
	ID     uuid.UUID `json:"id"`
	Name   *string   `json:"name"`
	IsNSFW *bool     `json:"is_nsfw"`
}

func (h *Handler) updateTag(ctx echo.Context) error {
	var reqInput updateTagRequest
	err := ctx.Bind(&reqInput)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	reqInput.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	input, err := models.NewUpdateTagInput(reqInput.ID, reqInput.Name, reqInput.IsNSFW)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	err = h.services.Tag.Update(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTagNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}
	ctx.NoContent(http.StatusOK)
	return nil
}
func (h *Handler) deleteTag(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	err = h.services.Tag.Delete(ctx.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTagNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}
	ctx.NoContent(http.StatusOK)
	return nil
}
