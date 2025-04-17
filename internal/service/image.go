package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"image"
	"io"
)

type ImageService struct {
	ImageRepository domain.ImageRepository
	S3              *minio.Client
}

func (s *ImageService) FindById(ctx context.Context, id string) (*domain.Image, error) {

	image, err := s.ImageRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Println(image.Path)

	obj, err := s.S3.GetObject(ctx, "notsamsa", image.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	abcde, _ := obj.Stat()
	fmt.Println(abcde.ContentType)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, obj)
	if err != nil {
		return nil, err
	}

	fmt.Printf("LEn: %d", buf.Len())

	image.ImageBinary = &buf
	return image, nil
}

func (s *ImageService) Create(ctx context.Context, image *domain.Image) error {
	image.Id = uuid.New().String()

	image.Path = fmt.Sprintf("collection/%s/image/%s", image.ImageCollectionId, image.Id)

	_, err := s.S3.PutObject(
		ctx,
		"notsamsa",
		image.Path,
		image.ImageBinary,
		image.Size,
		minio.PutObjectOptions{
			ContentType: fmt.Sprintf("image/%s", image.Format),
		},
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.ImageRepository.Create(ctx, image)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageService) Delete(ctx context.Context, id string) error {

	im, err := s.ImageRepository.FindById(ctx, id)

	if err != nil {
		return err
	}

	err = s.ImageRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.S3.RemoveObject(ctx, "notsamsa", im.Path, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageService) FindByCollection(ctx context.Context, collectionId string) ([]*domain.Image, error) {
	return s.ImageRepository.FindByCollection(ctx, collectionId)
}

func (s *ImageService) decodeImage(reader io.Reader) (*image.Image, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &im, nil
}
