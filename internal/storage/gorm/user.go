package gormstore

import (
	"context"
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {

	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {

	var entity domain.User

	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity, nil

}

func (r *UserRepository) Create(ctx context.Context, entity *domain.User) error {

	tx := r.db.WithContext(ctx).Create(entity)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	var entity domain.User
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepository) Update(ctx context.Context, entity *domain.User) error {

	tx := r.db.WithContext(ctx).Save(entity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

func (r *UserRepository) Delete(ctx context.Context, id string) error {

	var entity domain.User
	tx := r.db.WithContext(ctx).Delete(&entity, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func (r *UserRepository) FindByMail(ctx context.Context, mail string) (*domain.User, error) {

	var user domain.User
	err := r.db.WithContext(ctx).Where("mail = ?", mail).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
