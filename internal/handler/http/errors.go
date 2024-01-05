package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
)

const (
	validationErrorType = "validationError"
	appErrorType        = "appError"
	authErrorType       = "authorizationError"
	baseErrorType       = "baseError"
)

type httpError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type errorResponse struct {
	Error httpError `json:"error"`
}

func (h *Handler) newValidationErrorResponse(ctx echo.Context, code int, err error) error {
	return ctx.JSON(code, errorResponse{
		Error: httpError{
			Message: err.Error(),
			Code:    code,
			Type:    validationErrorType,
		},
	})
}

func (h *Handler) newErrorResponse(ctx echo.Context, code int, err string) error {
	return ctx.JSON(code, errorResponse{
		Error: httpError{
			Message: err,
			Code:    code,
			Type:    baseErrorType,
		},
	})
}

func (h *Handler) newAppErrorResponse(ctx echo.Context, err error) error {
	switch {
	case errors.As(err, &apperror.DBError{}):
		dbErr := err.(apperror.DBError)
		h.logger.Error(dbErr.Err.Error(), map[string]interface{}{
			"service": dbErr.Service,
			"func":    dbErr.Func,
			"query":   dbErr.Query,
			"args":    dbErr.Args,
		})
	default:
		h.logger.Error(err.Error(), nil)
	}

	return ctx.JSON(http.StatusInternalServerError, errorResponse{
		Error: httpError{
			Message: "server error",
			Code:    http.StatusInternalServerError,
			Type:    appErrorType,
		},
	})
}

func (h *Handler) newAuthErrorResponse(ctx echo.Context, code int, err error) error {
	return ctx.JSON(code, errorResponse{
		Error: httpError{
			Message: err.Error(),
			Code:    code,
			Type:    authErrorType,
		},
	})
}
