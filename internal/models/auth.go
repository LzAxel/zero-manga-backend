package models

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
)

type LoginUserInput struct {
	Username string
	Password string
}

func NewLoginUserInput(username string, password string) (LoginUserInput, error) {
	input := LoginUserInput{
		Username: username,
		Password: password,
	}

	if err := validation.ValidateStruct(&input,
		validation.Field(&input.Username, validation.Required),
		validation.Field(&input.Password, validation.Required),
	); err != nil {
		return input, err
	}

	return input, nil
}
