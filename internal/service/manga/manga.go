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

type Uploader interface {
	UploadMangaPreview(ctx context.Context, mangaID guuid.UUID, file models.UploadFile) (filestorage.FileInfo, error)
	DeleteMangaPreview(ctx context.Context, mangaID guuid.UUID, previewID guuid.UUID) error
	DeleteManga(ctx context.Context, mangaID guuid.UUID) error
}

type Manga struct {
	repo         repository.Manga
	chaptersRepo repository.Chapter
	uploader     Uploader
}

func New(ctx context.Context, repository repository.Manga, chaptersRepo repository.Chapter, uploader Uploader) *Manga {
	return &Manga{
		repo:         repository,
		chaptersRepo: chaptersRepo,
		uploader:     uploader,
	}
}

func (m *Manga) Create(ctx context.Context, manga models.CreateMangaInput) error {
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
	}

	preview, err := m.uploader.UploadMangaPreview(ctx, dto.ID, manga.PreviewFile)
	if err != nil {
		return err
	}

	dto.PreviewURL = preview.URL

	err = m.repo.Create(ctx, dto)
	if err != nil {
		delErr := m.uploader.DeleteMangaPreview(ctx, dto.ID, preview.ID)
		if delErr != nil {
			return apperror.NewAppError(
				fmt.Errorf("failed to delete file on create manga error: %w: %w", delErr, err),
				"Manga",
				"Create",
				nil,
			)
		}

		return err
	}

	return err
}

func (m *Manga) Update(ctx context.Context, manga models.UpdateMangaInput) error {
	var (
		previewURL *string
		newSlug    *string
		err        error
	)

	if manga.Title != nil {
		newSlugVal := slug.GenerateSlug(*manga.Title)
		newSlug = &newSlugVal
	}
	var previewID guuid.UUID
	if manga.PreviewFile != nil {
		newPreview, err := m.uploader.UploadMangaPreview(ctx, manga.ID, *manga.PreviewFile)
		if err != nil {
			return err
		}
		previewID = newPreview.ID
		previewURL = &newPreview.URL
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
		PreviewURL:     previewURL,
	}

	err = m.repo.Update(ctx, dto)
	if err != nil {
		if manga.PreviewFile != nil {
			delErr := m.uploader.DeleteMangaPreview(ctx, dto.ID, previewID)
			if delErr != nil {
				return apperror.NewAppError(
					fmt.Errorf("failed to delete file on update manga error: %w: %w", delErr, err),
					"Manga",
					"Update",
					nil,
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

	chaptersCount, err := m.chaptersRepo.CountByManga(ctx, manga.ID)
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
		PreviewURL:     manga.PreviewURL,
		ChaptersCount:  chaptersCount,
	}, nil
}
func (m *Manga) GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.MangaOutput, uint64, error) {
	mangaList, count, err := m.repo.GetAll(ctx, pagination, filters)
	if err != nil {
		return nil, 0, err
	}

	newManga := make([]models.MangaOutput, len(mangaList))
	for i, manga := range mangaList {
		chaptersCount, err := m.chaptersRepo.CountByManga(ctx, manga.ID)
		if err != nil {
			return []models.MangaOutput{}, 0, err
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
			PreviewURL:     manga.PreviewURL,
			ChaptersCount:  chaptersCount,
		}
	}

	return newManga, count, err
}

func (m *Manga) Delete(ctx context.Context, id guuid.UUID) error {
	err := m.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("Manga.Delete: delete from db %w", err)
	}
	err = m.uploader.DeleteManga(ctx, id)
	if err != nil {
		return fmt.Errorf("Manga.Delete: delete from file storage %w", err)
	}

	return nil
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrMangaNotFound
	}
	return err
}
