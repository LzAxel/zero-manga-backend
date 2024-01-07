package models

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ChapterID uuid.UUID `db:"chapter_id" json:"chapter_id"`
	URL       string    `db:"url" json:"url"`
	Number    int       `db:"number" json:"number"`
	Height    int       `db:"height" json:"height"`
	Width     int       `db:"width" json:"width"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
