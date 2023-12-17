package models

import (
	"time"

	"github.com/google/uuid"
)

type Chapter struct {
	ID         uuid.UUID
	MangaID    uuid.UUID
	Title      *string
	Number     uint
	Volume     uint
	PageConut  uint
	FilePath   string
	UploaderID uuid.UUID
	UploadedAt time.Time
}
