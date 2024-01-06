package manga

import (
	"context"
	"errors"
	"fmt"

	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/slug"
	"github.com/lzaxel/zero-manga-backend/pkg/uuid"

	guuid "github.com/google/uuid"
)

type Manga struct {
	repo        repository.Manga
	fileStorage filestorage.FileStorage
}

func New(ctx context.Context, repository repository.Manga, fileStorage filestorage.FileStorage) *Manga {
	return &Manga{
		repo:        repository,
		fileStorage: fileStorage,
	}
}

func (m *Manga) Create(ctx context.Context, manga models.CreateMangaInput) error {
	fileID, err := m.fileStorage.SaveFile(filestorage.MangaBucket, manga.PreviewFile.Filename, manga.PreviewFile.Data)
	if err != nil {
		return err
	}

	dto := models.Manga{
		ID:             uuid.New(),
		Title:          manga.Title,
		SecondaryTitle: manga.SecondaryTitle,
		Description:    manga.Description,
		Slug:           slug.GenerateSlug(manga.Title),
		Type:           manga.Type,
		Status:         manga.Status,
		AgeRestrict:    manga.AgeRestrict,
		ReleaseYear:    manga.ReleaseYear,
		PreviewFileID:  fileID,
	}

	err = m.repo.Create(ctx, dto)
	if err != nil {
		delErr := m.fileStorage.DeleteFile(filestorage.MangaBucket, fileID)
		if delErr != nil {
			return apperror.NewAppError(
				fmt.Errorf("failed to delete file on create manga error: %w: %w", delErr, err),
				"Manga",
				"Create",
				map[string]any{"file_id": fileID.String()},
			)
		}

		return err
	}

	return err
}

func (m *Manga) Update(ctx context.Context, manga models.UpdateMangaInput) error {
	var (
		fileID  *guuid.UUID
		newSlug *string
	)

	if manga.PreviewFile != nil {
		previewFileID, err := m.fileStorage.SaveFile(filestorage.MangaBucket, manga.PreviewFile.Filename, manga.PreviewFile.Data)
		if err != nil {
			return err
		}
		fileID = &previewFileID
	}

	if manga.Title != nil {
		newSlugVal := slug.GenerateSlug(*manga.Title)
		newSlug = &newSlugVal
	}

	dto := models.UpdateMangaRecord{
		ID:             manga.ID,
		Title:          manga.Title,
		SecondaryTitle: manga.SecondaryTitle,
		Description:    manga.Description,
		Slug:           newSlug,
		Type:           manga.Type,
		Status:         manga.Status,
		AgeRestrict:    manga.AgeRestrict,
		ReleaseYear:    manga.ReleaseYear,
		PreviewFileID:  fileID,
	}

	err := m.repo.Update(ctx, dto)
	if err != nil {
		if fileID != nil {
			delErr := m.fileStorage.DeleteFile(filestorage.MangaBucket, *fileID)
			if delErr != nil {
				return apperror.NewAppError(
					fmt.Errorf("failed to delete file on update manga error: %w: %w", delErr, err),
					"Manga",
					"Update",
					map[string]any{"file_id": fileID.String()},
				)
			}
		}

		return err
	}

	return err
}

func (m *Manga) GetOne(ctx context.Context, filters models.MangaFilters) (models.MangaOutput, error) {
	manga, err := m.repo.GetOne(ctx, filters)
	if err != nil {
		return models.MangaOutput{}, handleNotFoundError(err)
	}
	previewURL, err := m.fileStorage.GetFileURL(filestorage.MangaBucket, manga.PreviewFileID)
	if err != nil {
		return models.MangaOutput{}, err
	}

	return models.MangaOutput{
		ID:             manga.ID,
		Title:          manga.Title,
		SecondaryTitle: manga.SecondaryTitle,
		Description:    manga.Description,
		Slug:           manga.Slug,
		Type:           manga.Type,
		Status:         manga.Status,
		AgeRestrict:    manga.AgeRestrict,
		ReleaseYear:    manga.ReleaseYear,
		PreviewURL:     previewURL,
	}, nil
}
func (m *Manga) GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.MangaOutput, uint64, error) {
	manga, count, err := m.repo.GetAll(ctx, pagination, filters)
	if err != nil {
		return nil, 0, err
	}

	newManga := make([]models.MangaOutput, len(manga))
	for i, manga := range manga {
		previewURL, err := m.fileStorage.GetFileURL(filestorage.MangaBucket, manga.PreviewFileID)
		if err != nil {
			return nil, 0, err
		}

		newManga[i] = models.MangaOutput{
			ID:             manga.ID,
			Title:          manga.Title,
			SecondaryTitle: manga.SecondaryTitle,
			Description:    manga.Description,
			Slug:           manga.Slug,
			Type:           manga.Type,
			Status:         manga.Status,
			AgeRestrict:    manga.AgeRestrict,
			ReleaseYear:    manga.ReleaseYear,
			PreviewURL:     previewURL,
		}
	}

	return newManga, count, handleNotFoundError(err)
}

func (m *Manga) Delete(ctx context.Context, id guuid.UUID) error {
	return m.repo.Delete(ctx, id)
}

func handleNotFoundError(err error) error {
	if err != nil {
		if errors.As(err, &apperror.DBError{}) {
			dbErr := err.(apperror.DBError)
			switch {
			case errors.Is(dbErr.Err, apperror.ErrNotFound):
				return models.ErrMangaNotFound
			default:
				return err
			}
		}
	}

	return nil
}
