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
	"net/url"
	"strings"
	"time"
)

type ImageService struct {
	ImageRepository domain.ImageRepository
	S3              *minio.Client
}

func (s *ImageService) FindById(ctx context.Context, id string, preview bool) (*domain.Image, error) {

	im, err := s.ImageRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	path := im.Path
	if preview {
		path = strings.Replace(path, "images", "previews", 1)
	}

	obj, err := s.S3.GetObject(ctx, "notsamsa", path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, obj)
	if err != nil {
		return nil, err
	}

	im.ImageBinary = &buf

	s.populateWithPresignedURLs(ctx, im)

	return im, nil
}

func (s *ImageService) Create(ctx context.Context, image *domain.Image) error {
	image.Id = uuid.New().String()

	image.Path = fmt.Sprintf("collection/%s/images/%s", image.CollectionId, image.Id)

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

	err = s.S3.RemoveObject(ctx, "notsamsa", strings.Replace(im.Path, "images", "previews", 1), minio.RemoveObjectOptions{})
	return nil
}

func (s *ImageService) FindByCollection(ctx context.Context, collectionId string) ([]*domain.Image, error) {
	images, err := s.ImageRepository.FindByCollection(ctx, collectionId)
	s.populateWithPresignedURLs(ctx, images...)
	return images, err
}

func (s *ImageService) decodeImage(reader io.Reader) (*image.Image, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &im, nil
}

func (s *ImageService) populateWithPresignedURLs(ctx context.Context, images ...*domain.Image) {
	for _, image := range images {
		url, err := s.S3.PresignedGetObject(
			ctx,
			"notsamsa",
			fmt.Sprintf("collection/%s/previews/%s", image.CollectionId, image.Id),
			time.Minute*5,
			make(url.Values),
		)

		if err != nil {
			fmt.Println(err)
		}

		urlString := url.String()
		fmt.Println(urlString)
		image.PreviewUrl = urlString
	}
}
