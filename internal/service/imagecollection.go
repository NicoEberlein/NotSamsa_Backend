package service

import (
	"context"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
)

type ImageCollectionService struct {
	ImageCollectionRepository domain.ImageCollectionRepository
}

func (s *ImageCollectionService) FindById(ctx context.Context, id string) (*domain.ImageCollection, error) {
	return s.ImageCollectionRepository.FindById(ctx, id)
}

func (s *ImageCollectionService) Create(ctx context.Context, collection *domain.ImageCollection) (string, error) {
	collection.Id = uuid.New().String()
	if err := s.ImageCollectionRepository.Create(ctx, collection); err != nil {
		return "", err
	} else {
		return collection.Id, nil
	}
}

func (s *ImageCollectionService) Update(ctx context.Context, collection *domain.ImageCollection) error {
	return s.ImageCollectionRepository.Update(ctx, collection)
}

func (s *ImageCollectionService) Delete(ctx context.Context, id string) error {
	return s.ImageCollectionRepository.Delete(ctx, id)
}

func (s *ImageCollectionService) FindByUser(ctx context.Context, userId string) ([]*domain.ImageCollection, error) {
	return s.ImageCollectionRepository.FindByUser(ctx, userId)
}
