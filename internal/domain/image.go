package domain

import (
	"context"
	"image"
	"time"
)

type Image struct {
	Id                string `gorm:"type:uuid;primary_key"`
	ImageCollectionId string
	ImageCollection   ImageCollection `gorm:"foreignKey:ImageCollectionId"`
	Path              string
	Filename          string
	Size              uint64
	Date              time.Time
	ImageBinary       *image.Image `gorm:"-"`
}

type ImageRepository interface {
	FindById(ctx context.Context, id string) (*Image, error)
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, image *Image) error
	Update(ctx context.Context, image *Image) error
	Delete(ctx context.Context, id string) error
}
