package filestorage

import (
	"errors"

	"github.com/google/uuid"
)

const (
	MangaBucket = "manga"
)

var (
	ErrNotFound = errors.New("file not found")
)

type FileStorage interface {
	SaveFile(bucket string, filename string, data []byte) (uuid.UUID, error)
	GetFile(bucket string, id uuid.UUID) ([]byte, error)
	GetFilePath(bucket string, id uuid.UUID) (string, error)
	GetFileURL(bucket string, id uuid.UUID) (string, error)
	DeleteFile(bucket string, id uuid.UUID) error
}
