package uploader

import (
	"context"
	"fmt"
	"path"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/models"
)

func (u *Uploader) UploadPage(
	ctx context.Context,
	mangaID uuid.UUID,
	chapterID uuid.UUID,
	file models.UploadFile) (filestorage.FileInfo, error) {

	pageFileInfo, err := u.fileStorage.SaveFile(u.formatMangaChapterBucket(mangaID, chapterID), file.Filename, file.Data)
	if err != nil {
		return filestorage.FileInfo{}, fmt.Errorf("Uploader.UploadPage: %w", err)
	}

	return pageFileInfo, nil
}

func (u *Uploader) formatMangaChapterBucket(mangaID uuid.UUID, chapterID uuid.UUID) string {
	return path.Join(u.formatMangaBucket(mangaID), u.formatChapterBucket(chapterID))
}

func (u *Uploader) formatChapterBucket(chapterID uuid.UUID) string {
	return path.Join(ChapterBucket, chapterID.String())
}
