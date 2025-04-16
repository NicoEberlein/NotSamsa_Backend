package gormstore

import (
	"context"
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"gorm.io/gorm"
)

type ImageCollectionRepository struct {
	db *gorm.DB
}

func NewImageCollectionRepository(db *gorm.DB) domain.ImageCollectionRepository {
	return &ImageCollectionRepository{db: db}
}

func (r *ImageCollectionRepository) FindById(ctx context.Context, id string) (*domain.ImageCollection, error) {

	var entity domain.ImageCollection

	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity, nil

}

func (r *ImageCollectionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var entity domain.ImageCollection
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ImageCollectionRepository) Create(ctx context.Context, entity *domain.ImageCollection) error {

	tx := r.db.WithContext(ctx).Create(entity)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (r *ImageCollectionRepository) Update(ctx context.Context, entity *domain.ImageCollection) error {

	tx := r.db.WithContext(ctx).Save(entity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

func (r *ImageCollectionRepository) Delete(ctx context.Context, id string) error {

	var entity domain.ImageCollection
	tx := r.db.WithContext(ctx).Delete(&entity, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func (r *ImageCollectionRepository) FindByUser(ctx context.Context, userId string) ([]*domain.ImageCollection, error) {

	imageCollections := make([]*domain.ImageCollection, 0)

	tx := r.db.WithContext(ctx).
		Model(&domain.ImageCollection{}).
		Preload("Owner").
		Where("owner_id = ?", userId).
		Find(&imageCollections)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return imageCollections, nil
}
