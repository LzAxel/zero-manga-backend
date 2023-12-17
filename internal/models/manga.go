package models

import "github.com/google/uuid"

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
	ID             uuid.UUID
	Title          string
	SecondaryTitle string
	Description    string
	Slug           string
	Type           NovelType
	Status         MangaStatus
	AgeRestrict    AgeRestrict
	ReleaseYear    uint16
	PreviewID      uuid.UUID
}
