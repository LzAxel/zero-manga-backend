package tag

import (
	"context"
	"errors"
	"fmt"

	"github.com/lzaxel/zero-manga-backend/internal/apperror"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
	"github.com/lzaxel/zero-manga-backend/pkg/slug"
	"github.com/lzaxel/zero-manga-backend/pkg/uuid"

	guuid "github.com/google/uuid"
)

type Tag struct {
	repo repository.Tag
}

func New(repo repository.Tag) *Tag {
	return &Tag{
		repo: repo,
	}
}

func (t *Tag) Create(ctx context.Context, tag models.CreateTagInput) error {
	dto := models.Tag{
		ID:        uuid.New(),
		Name:      tag.Name,
		Slug:      slug.GenerateSlug(tag.Name),
		IsNSFW:    tag.IsNSFW,
		CreatedAt: clock.Now(),
	}
	err := t.repo.Create(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTagDuplicated):
			return models.ErrTagDuplicated
		default:
			return apperror.NewAppError(
				fmt.Errorf("failed to create tag: %w", err),
				"Tag",
				"Create",
				nil,
			)
		}
	}
	return err
}

func (t *Tag) Update(ctx context.Context, tag models.UpdateTagInput) error {
	var newSlug *string
	if tag.Name != nil {
		newSlugVal := slug.GenerateSlug(*tag.Name)
		newSlug = &newSlugVal
	}
	dto := models.UpdateTagRecord{
		ID:     tag.ID,
		Name:   tag.Name,
		Slug:   newSlug,
		IsNSFW: tag.IsNSFW,
	}
	err := t.repo.Update(ctx, dto)
	if err != nil {
		return apperror.NewAppError(
			fmt.Errorf("failed to update tag: %w", err),
			"Tag",
			"Update",
			nil,
		)
	}
	return err
}
func (t *Tag) GetAll(ctx context.Context) ([]models.Tag, error) {
	tags, err := t.repo.GetAll(ctx)
	if err != nil {
		return nil, apperror.NewAppError(
			fmt.Errorf("failed to get tags: %w", err),
			"Tag",
			"GetAll",
			nil,
		)
	}
	return tags, err
}
func (t *Tag) Delete(ctx context.Context, id guuid.UUID) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		return apperror.NewAppError(
			fmt.Errorf("failed to delete tag: %w", err),
			"Tag",
			"Delete",
			map[string]any{"tag_id": id},
		)
	}
	return err

}
