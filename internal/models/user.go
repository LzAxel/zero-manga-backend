package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	usernameMinLength = 4
	passwordMinLength = 8
)

var (
	ErrInvalidGender   = errors.New("invalid gender")
	ErrPasswordShort   = errors.New("password is too short")
	ErrPasswordLong    = errors.New("password is too long")
	ErrInvalidPassword = errors.New("password contains invalid characters")
	ErrUsernameShort   = errors.New("username is too short")
	ErrUsernameLong    = errors.New("username is too long")
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

type User struct {
	ID           uuid.UUID
	Username     string
	DisplayName  string
	Bio          string
	Email        string
	Gender       GenderType
	Type         UserType
	AvatarID     uuid.UUID
	PasswordHash []byte
	OnlineAt     time.Time
	RegisteredAt time.Time
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
	gender GenderType,
	bio *string,
) (CreateUserInput, error) {
	if gender != GenderTypeMale && gender != GenderTypeFemale {
		return CreateUserInput{}, ErrInvalidGender
	}
	if len(username) < usernameMinLength {
		return CreateUserInput{}, ErrUsernameShort
	}
	if len(password) < passwordMinLength {
		return CreateUserInput{}, ErrPasswordShort
	}
	if strings.ContainsAny(password, "\"'()+,-./:;<=>?[\\]_`{|}~") {
		return CreateUserInput{}, ErrInvalidPassword
	}
	return CreateUserInput{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Password:    password,
		Gender:      gender,
		Bio:         bio,
	}, nil
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
