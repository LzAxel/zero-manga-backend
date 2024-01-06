package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrMangaNotFound    = errors.New("Manga not found")
	ErrMangaTitleExists = errors.New("Manga title already taken")
)

type NovelType uint8

const (
	MangaType = iota + 1
	ManhwaType
	ManhuaType
)

type AgeRestrict uint8

const (
	AgeRestrictNo = iota + 1
	AgeRestrict16
	AgeRestrict18
)

type MangaStatus uint8

const (
	StatusAnnounced = iota + 1
	StatusOngoing
	StatusPaused
	StatusStopped
	StatusCompleted
)

type Manga struct {
	ID             uuid.UUID   `db:"id"`
	Title          string      `db:"title"`
	SecondaryTitle *string     `db:"secondary_title"`
	Description    string      `db:"description"`
	Slug           string      `db:"slug"`
	Type           NovelType   `db:"type"`
	Status         MangaStatus `db:"status"`
	AgeRestrict    AgeRestrict `db:"age_restrict"`
	ReleaseYear    uint16      `db:"release_year"`
	PreviewFileID  uuid.UUID   `db:"preview_file_id"`
}

type MangaOutput struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	SecondaryTitle *string     `json:"secondary_title"`
	Description    string      `json:"description"`
	Slug           string      `json:"slug"`
	Type           NovelType   `json:"type"`
	Status         MangaStatus `json:"status"`
	AgeRestrict    AgeRestrict `json:"age_restrict"`
	ReleaseYear    uint16      `json:"release_year"`
	PreviewURL     string      `json:"preview_url"`
}

type UpdateMangaInput struct {
	ID             uuid.UUID
	Title          *string
	SecondaryTitle *string
	Description    *string
	Type           *NovelType
	Status         *MangaStatus
	AgeRestrict    *AgeRestrict
	ReleaseYear    *uint16
	PreviewFile    *UploadFile
}

type UpdateMangaRecord struct {
	ID             uuid.UUID
	Title          *string
	SecondaryTitle *string
	Description    *string
	Slug           *string
	Type           *NovelType
	Status         *MangaStatus
	AgeRestrict    *AgeRestrict
	ReleaseYear    *uint16
	PreviewFileID  *uuid.UUID
}

func NewUpdateMangaInput(
	id uuid.UUID,
	title *string,
	secondaryTitle *string,
	description *string,
	type_ *NovelType,
	status *MangaStatus,
	ageRestrict *AgeRestrict,
	releaseYear *uint16,
	previewFile *UploadFile,
) UpdateMangaInput {
	return UpdateMangaInput{
		ID:             id,
		Title:          title,
		SecondaryTitle: secondaryTitle,
		Description:    description,
		Type:           type_,
		Status:         status,
		AgeRestrict:    ageRestrict,
		ReleaseYear:    releaseYear,
		PreviewFile:    previewFile,
	}
}

type CreateMangaInput struct {
	Title          string
	SecondaryTitle *string
	Description    string
	Type           NovelType
	Status         MangaStatus
	AgeRestrict    AgeRestrict
	ReleaseYear    uint16
	PreviewFile    UploadFile
}

func NewCreateMangaInput(
	title string,
	secondaryTitle *string,
	description string,
	type_ NovelType,
	status MangaStatus,
	ageRestrict AgeRestrict,
	releaseYear uint16,
	previewFile UploadFile,
) CreateMangaInput {
	return CreateMangaInput{
		Title:          title,
		SecondaryTitle: secondaryTitle,
		Description:    description,
		Type:           type_,
		Status:         status,
		AgeRestrict:    ageRestrict,
		ReleaseYear:    releaseYear,
		PreviewFile:    previewFile,
	}
}

type MangaFilters struct {
	ID    *uuid.UUID `query:"id"`
	Title *string    `query:"title"`
	Slug  *string    `query:"slug"`
}
type MangaGetAllFilters struct {
	Type        []NovelType   `query:"type"`
	Status      []MangaStatus `query:"status"`
	AgeRestrict []AgeRestrict `query:"age_restrict"`
	ReleaseYear *int          `query:"release_year"`
}
