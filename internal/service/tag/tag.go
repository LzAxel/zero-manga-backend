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
	repo                 repository.Tag
	mangaTagRelationRepo repository.MangaTagRelation
}

func New(repo repository.Tag, mangaTagRelationRepo repository.MangaTagRelation) *Tag {
	return &Tag{
		repo:                 repo,
		mangaTagRelationRepo: mangaTagRelationRepo,
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
				fmt.Errorf("creating tag: %w", err),
				"Tag",
				"Create",
				nil,
			)
		}
	}
	return err
}

func (t *Tag) AddTagToManga(ctx context.Context, mangaID, tagID guuid.UUID) error {
	err := t.mangaTagRelationRepo.Create(ctx, models.MangaTagRelation{
		MangaID: mangaID,
		TagID:   tagID,
	})

	if err != nil {
		switch {
		case errors.Is(err, models.ErrMangaOrTagNotFound):
			return models.ErrMangaOrTagNotFound
		default:
			return apperror.NewAppError(
				fmt.Errorf("adding tag to manga: %w", err),
				"Tag",
				"AddTagToManga",
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
			fmt.Errorf("updatig tag: %w", err),
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
			fmt.Errorf("getting all tags: %w", err),
			"Tag",
			"GetAll",
			nil,
		)
	}
	return tags, err
}

func (t *Tag) GetAllByMangaID(ctx context.Context, mangaID guuid.UUID) ([]models.Tag, error) {
	tagRelations, err := t.mangaTagRelationRepo.GetAllByMangaID(ctx, mangaID)
	if err != nil {
		return nil, apperror.NewAppError(
			fmt.Errorf("getting all manga tag relations: %w", err),
			"Tag",
			"GetAllByMangaID",
			map[string]any{"manga_id": mangaID},
		)
	}
	tags := make([]models.Tag, len(tagRelations))
	for i, relation := range tagRelations {
		tag, err := t.repo.GetByID(ctx, relation.TagID)
		if err != nil {
			return nil, apperror.NewAppError(
				fmt.Errorf("getting tag: %w", err),
				"Tag",
				"GetAllByMangaID",
				map[string]any{"manga_id": mangaID, "tag_id": relation.TagID},
			)
		}

		tags[i] = tag
	}

	return tags, err
}

func (t *Tag) Delete(ctx context.Context, id guuid.UUID) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		return apperror.NewAppError(
			fmt.Errorf("deleting tag: %w", err),
			"Tag",
			"Delete",
			map[string]any{"tag_id": id},
		)
	}
	return err

}

func (t *Tag) RemoveTagFromManga(ctx context.Context, mangaID, tagID guuid.UUID) error {
	err := t.mangaTagRelationRepo.Delete(ctx, models.MangaTagRelation{
		MangaID: mangaID,
		TagID:   tagID,
	})
	if err != nil {
		return apperror.NewAppError(
			fmt.Errorf("deleting tag from manga: %w", err),
			"Tag",
			"RemoveTagFromManga",
			map[string]any{"manga_id": mangaID, "tag_id": tagID},
		)
	}
	return err
}
