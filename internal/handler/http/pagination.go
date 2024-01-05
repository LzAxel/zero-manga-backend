package http

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

func (h *Handler) WithPagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)
			if err != nil {
				page = 1
			}
			pageLimit, err := strconv.ParseUint(c.QueryParam("page_limit"), 10, 64)
			if err != nil {
				pageLimit = models.DefaultLimit
			}

			pagination, err := models.NewPagination(page, pageLimit)
			if err != nil {
				return h.newAppErrorResponse(c, err)
			}
			c.Set("pagination", pagination)

			if err := next(c); err != nil {
				c.Error(err)
			}
			h.logger.Debug("pagination", map[string]interface{}{
				"page":       page,
				"pageLimit":  pageLimit,
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})

			return nil
		}
	}
}

func getPaginationFromContext(ctx echo.Context) (models.Pagination, error) {
	pagination, ok := ctx.Get("pagination").(models.Pagination)
	if !ok {
		pagination = models.Pagination{
			Page:      1,
			PageLimit: models.DefaultLimit,
		}
	}

	if pagination.PageLimit > models.MaxLimit {
		return models.Pagination{}, models.ErrInvalidLimit
	}

	return pagination, nil
}
