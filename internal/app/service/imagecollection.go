package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"net/url"
	"time"
)

type CollectionService struct {
	CollectionRepository domain.CollectionRepository
	UserRepository       domain.UserRepository
	ImageRepository      domain.ImageRepository
	S3                   *minio.Client
}

func (s *CollectionService) FindById(ctx context.Context, id string, userId *string) (*domain.Collection, error) {
	collection, err := s.CollectionRepository.FindById(ctx, id)
	s.populateWithPresignedURLs(ctx, collection)
	if userId != nil {
		s.setOwnerStatus(*userId, collection)
	}

	return collection, err
}

func (s *CollectionService) FindByUser(ctx context.Context, userId string) ([]*domain.Collection, error) {
	collections, err := s.CollectionRepository.FindByUser(ctx, userId)
	s.populateWithPresignedURLs(ctx, collections...)
	s.setOwnerStatus(userId, collections...)

	return collections, err
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

	collection, err := s.CollectionRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	for _, image := range collection.Images {
		err = s.ImageRepository.Delete(ctx, image.Id)
		if err != nil {
			return err
		}
	}

	return s.CollectionRepository.Delete(ctx, id)
}

func (s *CollectionService) AddParticipant(ctx context.Context, collectionId string, userMail string) error {

	collection, err := s.CollectionRepository.FindById(ctx, collectionId)
	if err != nil {
		return err
	}

	user, err := s.UserRepository.FindByMail(ctx, userMail)
	if err != nil {
		return err
	}

	if collection.OwnerId == user.Id {
		return errors.New("you can't add an account to his own collection")
	}

	return s.CollectionRepository.AddParticipant(ctx, collectionId, user.Id)
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

func (s *CollectionService) populateWithPresignedURLs(ctx context.Context, collections ...*domain.Collection) {
	for _, collection := range collections {
		if collection.PreviewImageId != nil {
			url, err := s.S3.PresignedGetObject(
				ctx,
				"notsamsa",
				fmt.Sprintf("collection/%s/previews/%s", collection.Id, *collection.PreviewImageId),
				time.Minute*5,
				make(url.Values),
			)

			if err != nil {
				continue
			}

			urlString := url.String()
			collection.PreviewImageUrl = &urlString
		}
	}
}

func (s *CollectionService) setOwnerStatus(userId string, collections ...*domain.Collection) {
	for _, collection := range collections {
		collection.OwnerStatus = collection.OwnerId == userId
	}
}
