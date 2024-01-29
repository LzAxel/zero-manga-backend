package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var (
	ErrTagExists   = errors.New("Tag already exitst")
	ErrTagNotFound = errors.New("Tag not found")
)

type Tag struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	IsNSFW    bool      `db:"is_nsfw" json:"is_nsfw"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UpdateTagRecord struct {
	ID     uuid.UUID
	Name   *string
	Slug   *string
	IsNSFW *bool
}

type CreateTagInput struct {
	Name   string
	IsNSFW bool
}

type UpdateTagInput struct {
	ID     uuid.UUID
	Name   *string
	IsNSFW *bool
}

func NewUpdateTagInput(
	ID uuid.UUID,
	name *string,
	isNSFW *bool,
) (UpdateTagInput, error) {
	input := UpdateTagInput{
		ID:     ID,
		Name:   name,
		IsNSFW: isNSFW,
	}
	return input, input.Validate()
}

func NewCreateTagInput(
	name string,
	isNSFW bool,
) (CreateTagInput, error) {
	input := CreateTagInput{
		Name:   name,
		IsNSFW: isNSFW,
	}
	return input, input.Validate()
}

func (input CreateTagInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Name,
			validation.Length(3, 30),
		))
}
func (input UpdateTagInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Name,
			validation.Length(3, 30),
		))
}
