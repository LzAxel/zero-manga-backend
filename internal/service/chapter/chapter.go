package chapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
	"github.com/lzaxel/zero-manga-backend/pkg/file"
)

type Uploader interface {
	UploadPage(ctx context.Context, mangaID uuid.UUID, chapterID uuid.UUID, file models.UploadFile) (filestorage.FileInfo, error)
}

type Chapter struct {
	repo      repository.Chapter
	pageRepo  repository.Page
	mangaRepo repository.Manga
	uploader  Uploader
}

func New(ctx context.Context,
	repository repository.Chapter,
	pageRepo repository.Page,
	mangaRepo repository.Manga,
	uploader Uploader,
) *Chapter {
	return &Chapter{
		repo:      repository,
		pageRepo:  pageRepo,
		mangaRepo: mangaRepo,
		uploader:  uploader,
	}
}

var (
	validPageExtensions = []string{".jpg", ".jpeg", ".png"}
)

func (c *Chapter) Create(ctx context.Context, chapter models.CreateChapterInput) error {
	_, err := c.mangaRepo.GetOne(ctx, models.MangaFilters{ID: &chapter.MangaID})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return models.ErrMangaNotFound
		}
		return fmt.Errorf("Chapter.Create: manga.GetOne: %w", err)
	}

	dto := models.Chapter{
		ID:         uuid.New(),
		MangaID:    chapter.MangaID,
		Title:      chapter.Title,
		Volume:     chapter.Volume,
		Number:     chapter.Number,
		UploaderID: chapter.UploaderID,
		UploadedAt: clock.Now(),
	}

	zip, err := file.GetFilesFromZip(chapter.PageArchive.Reader)
	if err != nil {
		return fmt.Errorf("Chapter.CreatePagesFromZip: file.GetFilesFromZip: %w", err)
	}

	// TODO: remove using 2 loop
	var pageCount uint = countValidImagesInZip(zip)
	if pageCount < 1 {
		return fmt.Errorf("Chapter.Create: %w", models.ErrNoValidImages)
	}
	dto.PageCount = pageCount

	file.SortZipFilesNumerically(zip)

	// TODO: add transaction
	err = c.repo.Create(ctx, dto)
	if err != nil {
		return fmt.Errorf("Chapter.Create: Create chapter: %w", err)
	}

	err = c.CreatePagesFromZip(ctx, CreatePagesFromZipInput{
		ChapterID: dto.ID,
		MangaID:   dto.MangaID,
		Files:     zip,
	})
	if err != nil {
		return fmt.Errorf("Chapter.Create: %w", err)
	}

	return nil
}

func (c *Chapter) GetAllByManga(ctx context.Context, pagination models.DBPagination, mangaID uuid.UUID) ([]models.Chapter, uint64, error) {
	_, err := c.mangaRepo.GetOne(ctx, models.MangaFilters{ID: &mangaID})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, 0, models.ErrMangaNotFound
		}
	}
	mangaList, count, err := c.repo.GetAllByManga(ctx, pagination, mangaID)
	if err != nil {
		return nil, 0, err
	}

	return mangaList, count, nil
}

func (c *Chapter) Get(ctx context.Context, chapterID uuid.UUID) (models.ChapterOutput, error) {
	chapter, err := c.repo.GetByID(ctx, chapterID)
	if err != nil {
		return models.ChapterOutput{}, fmt.Errorf("Chapter.Get: GetByID: %w", handleNotFoundError(err))
	}
	pages, err := c.pageRepo.GetAllByChapter(ctx, chapterID)
	if err != nil {
		return models.ChapterOutput{}, fmt.Errorf("Chapter.Get: page.GetAllByChapter: %w", err)
	}
	return models.ChapterOutput{
		ID:         chapter.ID,
		MangaID:    chapter.MangaID,
		Title:      chapter.Title,
		Volume:     chapter.Volume,
		Number:     chapter.Number,
		UploadedAt: chapter.UploadedAt,
		PageCount:  chapter.PageCount,
		Pages:      pages,
	}, nil
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrChapterNotFound
	}
	return err
}
