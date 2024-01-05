package models

import (
	"errors"
	"time"

	val "github.com/lzaxel/zero-manga-backend/internal/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type UserType uint8

const (
	UserTypeReader = iota + 1
	UserTypeEditor
	UserTypeModerator
	UserTypeAdmin
)

type GenderType uint8

const (
	GenderTypeMale = iota + 1
	GenderTypeFemale
)

var (
	ErrUsernameEmailExists = errors.New("Username or email already taken")
)

type UserFilters struct {
	OnlineAt     *time.Time `query:"online_at"`
	RegisteredAt *time.Time `query:"registered_at"`
	Type         []uint8    `query:"type"`
	Gender       []uint8    `query:"gender"`
}

type GetUserOutput struct {
	ID           uuid.UUID  `db:"id"`
	Username     string     `db:"username"`
	DisplayName  string     `db:"display_name"`
	Bio          string     `db:"bio"`
	Gender       GenderType `db:"gender"`
	Type         UserType   `db:"type"`
	AvatarID     uuid.UUID  `db:"avatar_id"`
	OnlineAt     time.Time  `db:"online_at"`
	RegisteredAt time.Time  `db:"registered_at"`
}

type User struct {
	ID           uuid.UUID  `db:"id"`
	Username     string     `db:"username"`
	DisplayName  string     `db:"display_name"`
	Bio          string     `db:"bio"`
	Email        string     `db:"email"`
	Gender       GenderType `db:"gender"`
	Type         UserType   `db:"type"`
	AvatarID     uuid.UUID  `db:"avatar_id"`
	PasswordHash []byte     `db:"password_hash"`
	OnlineAt     time.Time  `db:"online_at"`
	RegisteredAt time.Time  `db:"registered_at"`
}

type CreateUserInput struct {
	Username    string
	DisplayName *string
	Email       string
	Password    string
	Gender      GenderType
	Bio         *string
}

func NewCreateUserInput(
	username string,
	displayName *string,
	email string,
	password string,
	gender uint8,
	bio *string,
) (CreateUserInput, error) {
	input := CreateUserInput{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Password:    password,
		Gender:      GenderType(gender),
		Bio:         bio,
	}

	return input, input.Validate()
}

func (input CreateUserInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Username, val.UsernameRules...),
		validation.Field(&input.DisplayName, validation.Length(0, 32)),
		validation.Field(&input.Email, val.EmailRules...),
		validation.Field(&input.Password, val.PasswordRules...),
		validation.Field(&input.Gender, validation.In(
			GenderType(GenderTypeMale), GenderType(GenderTypeFemale),
		)),
		validation.Field(&input.Bio, validation.Length(0, 300)),
	)
}

type CreateUserRecord struct {
	ID           uuid.UUID
	Username     string
	DisplayName  *string
	Email        string
	PasswordHash []byte
	Gender       int
	Bio          *string
	Type         int
	OnlineAt     time.Time
	RegisteredAt time.Time
}

func NewCreateUserRecord(
	id uuid.UUID,
	username string,
	displayName *string,
	email string,
	passwordHash []byte,
	gender int,
	bio *string,
	userType int,
	onlineAt time.Time,
	registeredAt time.Time,
) CreateUserRecord {
	return CreateUserRecord{
		ID:           id,
		Username:     username,
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: passwordHash,
		Gender:       gender,
		Bio:          bio,
		Type:         userType,
		OnlineAt:     onlineAt,
		RegisteredAt: registeredAt,
	}
}
