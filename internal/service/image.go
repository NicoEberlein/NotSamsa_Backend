package service

import (
	"context"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"image"
)

type ImageService struct {
	ImageRepository domain.ImageRepository
	S3              *minio.Client
}

func (s *ImageService) FindById(ctx context.Context, id string) (*domain.Image, error) {
	return s.ImageRepository.FindById(ctx, id)
}

func (s *ImageService) Create(ctx context.Context, image *domain.Image) error {
	image.Id = uuid.New().String()
	return s.ImageRepository.Create(ctx, image)
}

func (s *ImageService) Update(ctx context.Context, image *domain.Image) error {
	return s.ImageRepository.Update(ctx, image)
}

func (s *ImageService) Delete(ctx context.Context, id string) error {
	return s.ImageRepository.Delete(ctx, id)
}

func (s *ImageService) uploadImage(ctx context.Context, image *image.Image) error {
	if err := s.S3.MakeBucket(ctx, "notsamsa", minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := s.S3.BucketExists(ctx, "notsamsa")
		if !(errBucketExists == nil && exists) {
			return errBucketExists
		}
	}
	return nil
}
