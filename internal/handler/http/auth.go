package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

type signUpRequest struct {
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name,omitempty"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	// Gender type
	// * 1 - Male
	// * 2 - Female
	Gender uint8   `json:"gender"`
	Bio    *string `json:"bio,omitempty"`
}

func (h *Handler) signUp(ctx echo.Context) error {
	var req signUpRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	input, err := models.NewCreateUserInput(
		req.Username,
		req.DisplayName,
		req.Email,
		req.Password,
		req.Gender,
		req.Bio,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = h.services.Authorization.Register(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameEmailExists):
			return h.newAuthErrorResponse(ctx, http.StatusConflict, err)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) signIn(ctx echo.Context) error {
	var req signInRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	input, err := models.NewLoginUserInput(
		req.Username,
		req.Password,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	tokenPair, err := h.services.Authorization.Login(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			return h.newAuthErrorResponse(ctx, http.StatusUnauthorized, err)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})

	return nil
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refreshTokens(ctx echo.Context) error {
	var req refreshRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	tokenPair, err := h.services.Authorization.RefreshTokens(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			return h.newAuthErrorResponse(ctx, http.StatusUnauthorized, models.ErrInvalidCredentials)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})

	return nil
}
