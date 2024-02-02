package manga

import (
	"context"
	"errors"
	"fmt"
	"slices"

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
	gradeRepo    repository.Grade
	tagRepo      repository.Tag
	uploader     Uploader
}

func New(
	repository repository.Manga,
	chaptersRepo repository.Chapter,
	gradeRepo repository.Grade,
	tagRepo repository.Tag,
	uploader Uploader,
) *Manga {
	return &Manga{
		repo:         repository,
		chaptersRepo: chaptersRepo,
		gradeRepo:    gradeRepo,
		uploader:     uploader,
		tagRepo:      tagRepo,
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
	err = m.repo.AddTagsToManga(ctx, dto.ID, manga.Tags...)
	if err != nil {
		return apperror.NewAppError(
			fmt.Errorf("adding tag to manga: %w", err),
			"Manga",
			"Create",
			map[string]interface{}{
				"manga_id": dto.ID,
				"tag_ids":  manga.Tags,
			},
		)
	}

	return nil
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

	if manga.Tags != nil {
		existedTags, err := m.tagRepo.GetAllByMangaID(ctx, manga.ID)
		if err != nil {
			return apperror.NewAppError(
				fmt.Errorf("getting manga tags: %w", err),
				"Manga",
				"Update",
				map[string]interface{}{
					"manga_id": manga.ID,
				},
			)
		}

		// removing tags if not in new tags list
		var toDeleteTags = make(guuid.UUIDs, 0)
		for _, existedTag := range existedTags {
			if !slices.Contains(*manga.Tags, existedTag.ID) {
				toDeleteTags = append(toDeleteTags, existedTag.ID)
			}
		}

		err = m.repo.RemoveTagsFromManga(ctx, manga.ID, toDeleteTags...)
		if err != nil {
			return apperror.NewAppError(
				fmt.Errorf("deleting tag: %w", err),
				"Manga",
				"Update",
				map[string]interface{}{
					"manga_id": manga.ID,
					"tag_ids":  toDeleteTags,
				},
			)
		}

		// adding new tags if not in existed tags
		var toAddTags = make(guuid.UUIDs, 0)
		for _, tagID := range *manga.Tags {
			var alreadyExists = false
			for _, existedTag := range existedTags {
				if tagID == existedTag.ID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				toAddTags = append(toAddTags, tagID)
			}
		}

		err = m.repo.AddTagsToManga(ctx, manga.ID, toAddTags...)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrMangaOrTagNotFound):
				return models.ErrTagNotFound
			default:
				return apperror.NewAppError(
					fmt.Errorf("adding tag to manga: %w", err),
					"Manga",
					"Update",
					map[string]interface{}{
						"manga_id": manga.ID,
						"tag_ids":  toAddTags,
					},
				)
			}
		}
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
	logMap := map[string]interface{}{
		"manga_id": filters.ID,
		"title":    filters.Title,
		"slug":     filters.Slug,
	}

	manga, err := m.repo.GetOne(ctx, filters)
	if err != nil {
		return models.MangaOutput{}, handleNotFoundError(err)
	}

	chaptersCount, err := m.chaptersRepo.CountByManga(ctx, manga.ID)
	if err != nil {
		return models.MangaOutput{}, apperror.NewAppError(
			fmt.Errorf("counting chapters: %w", err),
			"Manga",
			"GetOne",
			logMap,
		)
	}

	tags, err := m.tagRepo.GetAllByMangaID(ctx, manga.ID)
	if err != nil {
		return models.MangaOutput{}, apperror.NewAppError(
			fmt.Errorf("getting tags: %w", err),
			"Manga",
			"GetOne",
			logMap,
		)
	}

	avgGrade, gradeCount, err := m.gradeRepo.GetInfoByManga(ctx, manga.ID)
	if err != nil {
		return models.MangaOutput{}, apperror.NewAppError(
			fmt.Errorf("getting grade info: %w", err),
			"Manga",
			"GetOne",
			logMap,
		)
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
		Grade: models.GradeInfo{
			AvgGrade: avgGrade,
			Count:    gradeCount,
		},
		Tags: tags,
	}, nil
}
func (m *Manga) GetAll(ctx context.Context, pagination models.DBPagination, filters models.MangaGetAllFilters) ([]models.MangaOutput, uint64, error) {
	logMap := map[string]interface{}{
		"type":         filters.Type,
		"status":       filters.Status,
		"age_restrict": filters.AgeRestrict,
		"release_year": filters.ReleaseYear,
	}

	mangaList, count, err := m.repo.GetAll(ctx, pagination, filters)
	if err != nil {
		return nil, 0, apperror.NewAppError(
			fmt.Errorf("getting manga: %w", err),
			"Manga",
			"GetAll",
			logMap,
		)
	}

	newManga := make([]models.MangaOutput, len(mangaList))
	for i, manga := range mangaList {
		chaptersCount, err := m.chaptersRepo.CountByManga(ctx, manga.ID)

		if err != nil {
			return []models.MangaOutput{}, 0, apperror.NewAppError(
				fmt.Errorf("counting chapters: %w", err),
				"Manga",
				"GetAll",
				logMap,
			)
		}

		tags, err := m.tagRepo.GetAllByMangaID(ctx, manga.ID)
		if err != nil {
			return []models.MangaOutput{}, 0, apperror.NewAppError(
				fmt.Errorf("getting tags: %w", err),
				"Manga",
				"GetAll",
				logMap,
			)
		}

		avgGrade, gradeCount, err := m.gradeRepo.GetInfoByManga(ctx, manga.ID)
		if err != nil {
			return []models.MangaOutput{}, 0, apperror.NewAppError(
				fmt.Errorf("getting grade info: %w", err),
				"Manga",
				"GetAll",
				logMap,
			)
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
			Grade: models.GradeInfo{
				AvgGrade: avgGrade,
				Count:    gradeCount,
			},
			Tags: tags,
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
