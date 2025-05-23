package gormstore

import (
	"context"
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
)

type ImageCollectionRepository struct {
	db *gorm.DB
}

func NewImageCollectionRepository(db *gorm.DB) domain.CollectionRepository {
	return &ImageCollectionRepository{db: db}
}

func (r *ImageCollectionRepository) FindById(ctx context.Context, id string) (*domain.Collection, error) {

	var entity domain.Collection

	if err := r.db.
		WithContext(ctx).
		Preload(clause.Associations).
		First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}

	return &entity, nil

}

func (r *ImageCollectionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var entity domain.Collection
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ImageCollectionRepository) Create(ctx context.Context, entity *domain.Collection) error {

	tx := r.db.WithContext(ctx).Create(entity)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		} else if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		} else {
			return tx.Error
		}
	}

	return nil

}

func (r *ImageCollectionRepository) Update(ctx context.Context, entity *domain.Collection) error {

	tx := r.db.WithContext(ctx).Save(entity)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		} else {
			return tx.Error
		}
	}

	return nil

}

func (r *ImageCollectionRepository) Delete(ctx context.Context, id string) error {

	var entity domain.Collection
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

func (r *ImageCollectionRepository) FindByUser(ctx context.Context, userId string) ([]*domain.Collection, error) {

	imageCollectionsOwner := make([]*domain.Collection, 0)
	imageCollectionsParticipant := make([]*domain.Collection, 0)

	tx := r.db.WithContext(ctx).
		Model(&domain.Collection{}).
		Preload(clause.Associations).
		Where("owner_id = ?", userId).
		Find(&imageCollectionsOwner)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, tx.Error
		}
	}

	tx = r.db.WithContext(ctx).
		Model(&domain.Collection{}).
		Preload(clause.Associations).
		Joins("JOIN collection_participants ON collection_participants.collection_id = collections.id").
		Where("collection_participants.user_id = ?", userId).
		Find(&imageCollectionsParticipant)

	if tx.Error != nil {
		return nil, tx.Error
	}

	allCollections := slices.Concat(imageCollectionsOwner, imageCollectionsParticipant)

	return allCollections, nil
}

func (r *ImageCollectionRepository) AddParticipant(ctx context.Context, collectionId string, userId string) error {

	if err := r.db.WithContext(ctx).
		Model(&domain.Collection{Id: collectionId}).
		Association("Participants").
		Append(&domain.User{Id: userId}); err != nil {

		return err
	}

	return nil
}

func (r *ImageCollectionRepository) DeleteParticipant(ctx context.Context, collectionId string, userId string) error {
	collection, err := r.FindById(ctx, collectionId)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(collection).
		Association("Participants").
		Delete(&domain.User{Id: userId}); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		} else {
			return err
		}
	}

	return nil
}
