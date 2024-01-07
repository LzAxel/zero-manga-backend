package uploader

import (
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
)

const (
	PreviewBucket = "preview"
	MangaBucket   = "manga"
	ChapterBucket = "chapter"
)

type Uploader struct {
	fileStorage filestorage.FileStorage
}

func NewUploader(fileStorage filestorage.FileStorage) *Uploader {
	return &Uploader{
		fileStorage: fileStorage,
	}
}
