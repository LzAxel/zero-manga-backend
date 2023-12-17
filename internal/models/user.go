package models

import (
	"time"

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
