package domain

import (
	"context"
)

type Collection struct {
	Id              string   `gorm:"type:uuid;primary_key" json:"id"`
	OwnerId         string   `json:"-"`
	Owner           *User    `gorm:"foreignKey:OwnerId" json:"-"`
	Name            string   `json:"name"`
	Description     string   `gorm:"type:text" json:"description,omitempty"`
	Latitude        *float64 `gorm:"type:numeric(10,7)" json:"latitude,omitempty"`
	Longitude       *float64 `gorm:"type:numeric(10,7)" json:"longitude,omitempty"`
	Participants    []*User  `gorm:"many2many:collection_participants;" json:"-"`
	Images          []*Image `json:"-"`
	PreviewImageId  *string  `gorm:"type:uuid;column:preview_image_id" json:"-"`
	PreviewImageUrl *string  `gorm:"-" json:"previewImageUrl,omitempty"`
	OwnerStatus     bool     `gorm:"-" json:"owner"`
}

type CollectionRepository interface {
	FindById(ctx context.Context, id string) (*Collection, error)
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, collection *Collection) error
	Update(ctx context.Context, collection *Collection) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userId string) ([]*Collection, error)
	AddParticipant(ctx context.Context, collectionId string, userId string) error
	DeleteParticipant(ctx context.Context, collectionId string, userId string) error
}
