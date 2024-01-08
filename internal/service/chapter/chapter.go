package chapter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"

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
	zip, err := file.GetFilesFromZip(chapter.PageArchiveFile.Data)
	if err != nil {
		return fmt.Errorf("Chapter.Create: GetFilesFromZip: %w", err)
	}

	file.SortZipFilesNumerically(zip)

	dto := models.Chapter{
		ID:         uuid.New(),
		MangaID:    chapter.MangaID,
		Title:      chapter.Title,
		Volume:     chapter.Volume,
		Number:     chapter.Number,
		UploaderID: chapter.UploaderID,
		UploadedAt: clock.Now(),
	}

	// TODO: remove using 2 loop
	var pageCount uint
	for _, file := range zip {
		if isValidPageExtensions(file.Name) {
			pageCount++
		}
	}
	dto.PageCount = pageCount

	err = c.repo.Create(ctx, dto)
	if err != nil {
		return fmt.Errorf("Chapter.Create: Create: %w", err)
	}

	for i, file := range zip {
		if !isValidPageExtensions(file.Name) {
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("Chapter.Create: Open file: %w", err)
		}
		defer fileReader.Close()
		fileBytes, err := io.ReadAll(fileReader)
		if err != nil {
			return fmt.Errorf("Chapter.Create: Read file: %w", err)
		}

		uploadedPageFile, err := c.uploader.UploadPage(ctx, chapter.MangaID, dto.ID, models.UploadFile{
			Filename: file.Name,
			Data:     fileBytes,
		})
		if err != nil {
			return fmt.Errorf("Chapter.Create: UploadPage: %w", err)
		}

		err = c.pageRepo.Create(ctx, models.Page{
			ID:        uuid.New(),
			ChapterID: dto.ID,
			URL:       uploadedPageFile.URL,
			// TODO: add height and width
			Number:    i + 1,
			CreatedAt: clock.Now(),
		})
		if err != nil {
			return fmt.Errorf("Chapter.Create: Create page: %w", err)
		}
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

func isValidPageExtensions(filename string) bool {
	return slices.Contains(validPageExtensions, filepath.Ext(filename))
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrChapterNotFound
	}
	return err
}
