package domain

import (
	"bytes"
	"context"
	"time"
)

type Image struct {
	Id                string          `gorm:"type:uuid;primary_key"`
	ImageCollectionId string          `json:"-"`
	ImageCollection   ImageCollection `gorm:"foreignKey:ImageCollectionId" json:"-"`
	Path              string          `json:"-"`
	Name              string
	Size              int64
	Format            string
	UploadedAt        time.Time
	ImageBinary       *bytes.Buffer `gorm:"-" json:"-"`
}

func NewImage(imageCollectionId string, format string, name string, size int64, uploadedAt time.Time, imageBinary *bytes.Buffer) *Image {
	return &Image{
		ImageCollectionId: imageCollectionId,
		Format:            format,
		UploadedAt:        uploadedAt,
		Name:              name,
		Size:              size,
		ImageBinary:       imageBinary,
	}
}

type ImageRepository interface {
	FindById(ctx context.Context, id string) (*Image, error)
	Create(ctx context.Context, image *Image) error
	Update(ctx context.Context, image *Image) error
	Delete(ctx context.Context, id string) error
	FindByCollection(ctx context.Context, collectionId string) ([]*Image, error)
}
