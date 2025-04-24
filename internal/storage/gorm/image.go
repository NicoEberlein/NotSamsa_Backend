package gormstore

import (
	"context"
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) domain.ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

func (r *ImageRepository) FindById(ctx context.Context, id string) (*domain.Image, error) {

	var image *domain.Image

	if err := r.db.
		WithContext(ctx).
		Preload(clause.Associations).
		First(&image, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}

	return image, nil

}

func (r *ImageRepository) Create(ctx context.Context, image *domain.Image) error {

	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
	}

	return nil

}

func (r *ImageRepository) Update(ctx context.Context, image *domain.Image) error {

	tx := r.db.WithContext(ctx).Save(image)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		} else {
			return tx.Error
		}
	}

	return nil

}

func (r *ImageRepository) Delete(ctx context.Context, id string) error {

	var entity domain.Image
	tx := r.db.WithContext(ctx).Delete(&entity, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		} else {
			return tx.Error
		}
	}

	if tx.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil

}

func (r *ImageRepository) FindByCollection(ctx context.Context, collectionId string) ([]*domain.Image, error) {

	images := make([]*domain.Image, 0)

	tx := r.db.WithContext(ctx).
		Model(&domain.Image{}).
		Preload(clause.Associations).
		Where("collection_id = ?", collectionId).
		Find(&images)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, tx.Error
		}
	}

	return images, nil

}
