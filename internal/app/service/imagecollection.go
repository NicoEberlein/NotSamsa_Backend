package service

import (
	"context"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
)

type CollectionService struct {
	CollectionRepository domain.CollectionRepository
	UserRepository       domain.UserRepository
}

func (s *CollectionService) FindById(ctx context.Context, id string) (*domain.Collection, error) {
	return s.CollectionRepository.FindById(ctx, id)
}

func (s *CollectionService) Create(ctx context.Context, collection *domain.Collection) (string, error) {
	collection.Id = uuid.New().String()
	if err := s.CollectionRepository.Create(ctx, collection); err != nil {
		return "", err
	} else {
		return collection.Id, nil
	}
}

func (s *CollectionService) Patch(ctx context.Context, collection *domain.Collection) error {

	collectionFromDb, err := s.CollectionRepository.FindById(ctx, collection.Id)
	if err != nil {
		return err
	}

	if len(collection.Name) > 0 {
		collectionFromDb.Name = collection.Name
	}

	if len(collection.Description) > 0 {
		collectionFromDb.Description = collection.Description
	}

	if collection.Longitude != nil && collection.Latitude != nil {
		collectionFromDb.Longitude = collection.Longitude
		collectionFromDb.Latitude = collection.Latitude
	}

	if collection.PreviewImageId != nil {
		collectionFromDb.PreviewImageId = collection.PreviewImageId
	}

	return s.CollectionRepository.Update(ctx, collectionFromDb)
}

func (s *CollectionService) Delete(ctx context.Context, id string) error {
	return s.CollectionRepository.Delete(ctx, id)
}

func (s *CollectionService) FindByUser(ctx context.Context, userId string) ([]*domain.Collection, error) {
	return s.CollectionRepository.FindByUser(ctx, userId)
}

func (s *CollectionService) AddParticipant(ctx context.Context, collectionId string, userId string) error {

	if _, err := s.CollectionRepository.FindById(ctx, collectionId); err != nil {
		return err
	}

	if _, err := s.UserRepository.FindById(ctx, userId); err != nil {
		return err
	}

	return s.CollectionRepository.AddParticipant(ctx, collectionId, userId)
}

func (s *CollectionService) DeleteParticipant(ctx context.Context, collectionId string, userId string) error {

	if _, err := s.CollectionRepository.FindById(ctx, collectionId); err != nil {
		return err
	}

	if _, err := s.UserRepository.FindById(ctx, userId); err != nil {
		return err
	}

	return s.CollectionRepository.DeleteParticipant(ctx, collectionId, userId)
}
