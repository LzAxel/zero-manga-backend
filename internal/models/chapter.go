package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrChapterNotFound = errors.New("Chapter not found")
)

// TODO: create output struct with user name
type Chapter struct {
	ID         uuid.UUID `db:"id" json:"id"`
	MangaID    uuid.UUID `db:"manga_id" json:"manga_id"`
	Title      *string   `db:"title" json:"title"`
	Number     uint      `db:"number" json:"number"`
	Volume     uint      `db:"volume" json:"volume"`
	PageCount  uint      `db:"page_count" json:"page_count"`
	UploaderID uuid.UUID `db:"uploader_id" json:"uploader_id"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type ChapterOutput struct {
	ID         uuid.UUID `db:"id" json:"id"`
	MangaID    uuid.UUID `db:"manga_id" json:"manga_id"`
	Title      *string   `db:"title" json:"title"`
	Number     uint      `db:"number" json:"number"`
	Volume     uint      `db:"volume" json:"volume"`
	PageCount  uint      `db:"page_count" json:"page_count"`
	Pages      []Page    `db:"pages" json:"pages"`
	UploaderID uuid.UUID `db:"uploader_id" json:"uploader_id"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type ChapterFilter struct {
	MangaID uuid.UUID `query:"manga_id"`
	Number  uint      `query:"number"`
	Volume  uint      `query:"volume"`
}

type CreateChapterInput struct {
	MangaID         uuid.UUID
	UploaderID      uuid.UUID
	Title           *string
	Number          uint
	Volume          uint
	PageArchiveFile UploadFile // .zip file cantain pages
}
