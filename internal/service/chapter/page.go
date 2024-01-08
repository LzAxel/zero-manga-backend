package chapter

import (
	"archive/zip"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
)

type CreatePageInput struct {
	ChapterID uuid.UUID
	MangaID   uuid.UUID
	ImageFile models.UploadFile
	Number    int
}

func (c *Chapter) CreatePage(ctx context.Context, input CreatePageInput) error {
	uploadedPageFile, err := c.uploader.UploadPage(ctx, input.MangaID, input.ChapterID, models.UploadFile{
		Filename: input.ImageFile.Filename,
		Data:     input.ImageFile.Data,
	})
	if err != nil {
		return fmt.Errorf("Chapter.CreatePage: UploadPage: %w", err)
	}

	err = c.pageRepo.Create(ctx, models.Page{
		ID:        uuid.New(),
		ChapterID: input.ChapterID,
		URL:       uploadedPageFile.URL,
		// TODO: add height and width
		Number:    input.Number,
		CreatedAt: clock.Now(),
	})
	if err != nil {
		return fmt.Errorf("Chapter.CreatePage: Create page: %w", err)
	}

	return nil
}

type CreatePagesFromZipInput struct {
	ChapterID uuid.UUID
	MangaID   uuid.UUID
	Files     []*zip.File
}

func (c *Chapter) CreatePagesFromZip(ctx context.Context, input CreatePagesFromZipInput) error {
	// Add transaction
	for i, file := range input.Files {
		if !isValidPageExtensions(file.Name) {
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("Chapter.CreatePagesFromZip: Open file: %w", err)
		}
		defer fileReader.Close()
		fileBytes, err := io.ReadAll(fileReader)
		if err != nil {
			return fmt.Errorf("Chapter.CreatePagesFromZip: Read file: %w", err)
		}

		err = c.CreatePage(ctx, CreatePageInput{
			ChapterID: input.ChapterID,
			MangaID:   input.MangaID,
			ImageFile: models.UploadFile{
				Filename: file.Name,
				Data:     fileBytes,
			},
			Number: i + 1,
		})
		if err != nil {
			return fmt.Errorf("Chapter.CreatePagesFromZip: CreatePage: %w", err)
		}

	}

	return nil
}
