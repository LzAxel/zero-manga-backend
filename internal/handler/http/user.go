package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

type getAllUsersResponse struct {
	Users      []models.User         `json:"users"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getAllUsers(ctx echo.Context) error {
	var filters models.UserFilters

	err := ctx.Bind(&filters)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	h.logger.Debug("filters", map[string]interface{}{
		"filters": filters,
	})

	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	users, pagination, err := h.services.User.GetAll(ctx.Request().Context(), reqPagination, filters)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameEmailExists):
			return h.newAuthErrorResponse(ctx, http.StatusConflict, err)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getAllUsersResponse{Users: users, Pagination: pagination})

	return nil
}
