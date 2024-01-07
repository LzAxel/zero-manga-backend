package uploader

import (
	"context"
	"fmt"
	"path"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

const (
	PreviewBucket = "preview"
	MangaBucket   = "manga"
)

type Uploader struct {
	fileStorage filestorage.FileStorage
}

func NewUploader(fileStorage filestorage.FileStorage) *Uploader {
	return &Uploader{
		fileStorage: fileStorage,
	}
}

func (u *Uploader) UploadMangaPreview(ctx context.Context, mangaID uuid.UUID, file models.UploadFile) (filestorage.FileInfo, error) {
	previewFileInfo, err := u.fileStorage.SaveFile(u.formatMangaPreviewBucket(mangaID), file.Filename, file.Data)
	if err != nil {
		return filestorage.FileInfo{}, fmt.Errorf("Uploader.UploadMangaPreview: %w", err)
	}

	return previewFileInfo, nil
}

func (u *Uploader) DeleteMangaPreview(ctx context.Context, mangaID uuid.UUID, previewID uuid.UUID) error {
	err := u.fileStorage.DeleteFile(u.formatMangaPreviewBucket(mangaID), previewID)
	if err != nil {
		return fmt.Errorf("Uploader.DeleteMangaPreview: %w", err)
	}

	return nil
}

func (u *Uploader) DeleteManga(ctx context.Context, mangaID uuid.UUID) error {
	err := u.fileStorage.DeleteBucket(u.formatMangaBucket(mangaID))
	if err != nil {
		return fmt.Errorf("Uploader.DeleteManga: %w", err)
	}

	return nil
}

func (u *Uploader) formatMangaPreviewBucket(mangaID uuid.UUID) string {
	return path.Join(u.formatMangaBucket(mangaID), PreviewBucket)
}

func (u *Uploader) formatMangaBucket(mangaID uuid.UUID) string {
	return path.Join(MangaBucket, mangaID.String())
}
