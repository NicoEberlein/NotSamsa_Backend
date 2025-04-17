package domain

import (
	"context"
)

type ImageCollection struct {
	Id           string   `gorm:"type:uuid;primary_key" json:"id"`
	OwnerId      string   `json:"-"`
	Owner        *User    `gorm:"foreignKey:OwnerId" json:"-"`
	Name         string   `json:"name"`
	Participants []*User  `gorm:"many2many:image_participants;" json:"-"`
	Images       []*Image `json:"-"`
}

type ImageCollectionRepository interface {
	FindById(ctx context.Context, id string) (*ImageCollection, error)
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, collection *ImageCollection) error
	Update(ctx context.Context, collection *ImageCollection) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userId string) ([]*ImageCollection, error)
}
