package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type GradeType uint8

const (
	GradeTypeOne GradeType = iota + 1
	GradeTypeTwo
	GradeTypeThree
	GradeTypeFour
	GradeTypeFive
)

var (
	ErrInvalidGradeType  = errors.New("invalid grade type")
	ErrNotCreatorOfGrade = errors.New("not creator of grade")
	ErrDuplicatedGrade   = errors.New("grade duplicated")
	ErrGradeNotFound     = errors.New("grade not found")
)

type Grade struct {
	ID        int64
	UserID    uuid.UUID
	MangaID   uuid.UUID
	Grade     GradeType
	CreatedAt time.Time
}

func NewGrade(id int64, userID, mangaID uuid.UUID, gradeType GradeType, createdAt time.Time) Grade {
	return Grade{
		ID:        id,
		UserID:    userID,
		MangaID:   mangaID,
		Grade:     gradeType,
		CreatedAt: createdAt,
	}
}

type GradeInfo struct {
	AvgGrade float64 `json:"avg_grade"`
	Count    uint64  `json:"count"`
}

type CreateGradeInput struct {
	UserID  uuid.UUID
	MangaID uuid.UUID
	Grade   GradeType
}

func NewCreateGradeInput(userID uuid.UUID, mangaID string, grade uint8) (CreateGradeInput, error) {
	mangaUUID, err := uuid.Parse(mangaID)
	if err != nil {
		return CreateGradeInput{}, errors.New("invalid manga ID")
	}
	input := CreateGradeInput{
		UserID:  userID,
		MangaID: mangaUUID,
		Grade:   GradeType(grade),
	}
	return input, nil
}

type CreateGrade struct {
	UserID    uuid.UUID
	MangaID   uuid.UUID
	Grade     GradeType
	CreatedAt time.Time
}

func (i CreateGrade) Validate() error {
	if i.Grade < GradeTypeOne || i.Grade > GradeTypeFive {
		return ErrInvalidGradeType
	}
	return nil
}
func NewCreateGrade(userID, mangaID uuid.UUID, gradeType GradeType, createdAt time.Time) (CreateGrade, error) {
	grade := CreateGrade{
		UserID:    userID,
		MangaID:   mangaID,
		Grade:     gradeType,
		CreatedAt: createdAt,
	}

	return grade, grade.Validate()
}
