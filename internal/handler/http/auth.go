package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

type signUpRequest struct {
	Username    string  `json:"username" validate:"required"`
	DisplayName *string `json:"display_name,omitempty"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required"`
	// Gender type
	// * 1 - Male
	// * 2 - Female
	Gender uint8   `json:"gender" validate:"required" enums:"1,2"`
	Bio    *string `json:"bio,omitempty"`
}

// SignUp godoc
// @Summary      Create an account
// @Description  create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      201
// @Param        body    body      signUpRequest   true  "user data"
// @Router       /sign-up [post]
func (h *Handler) signUp(ctx echo.Context) error {
	var req signUpRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	input, err := models.NewCreateUserInput(
		req.Username,
		req.DisplayName,
		req.Email,
		req.Password,
		models.GenderType(req.Gender),
		req.Bio,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.services.User.Create(ctx.Request().Context(), input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}
