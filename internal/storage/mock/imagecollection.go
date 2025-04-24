package mock

import (
	"context"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
)

type CollectionRepository struct {
	collections map[string]*domain.Collection
}

func NewMockCollectionRepository() *CollectionRepository {
	return &CollectionRepository{
		collections: make(map[string]*domain.Collection),
	}
}

func (m *CollectionRepository) FindById(ctx context.Context, id string) (*domain.Collection, error) {
	collection, ok := m.collections[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return collection, nil
}

func (m *CollectionRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.collections[id]
	return ok, nil
}

func (m *CollectionRepository) Create(ctx context.Context, collection *domain.Collection) error {
	if _, ok := m.collections[collection.Id]; ok {
		return domain.ErrDuplicateEntity
	}
	m.collections[collection.Id] = collection
	return nil
}

func (m *CollectionRepository) Update(ctx context.Context, collection *domain.Collection) error {
	_, ok := m.collections[collection.Id]
	if !ok {
		return domain.ErrNotFound
	}
	m.collections[collection.Id] = collection
	return nil
}

func (m *CollectionRepository) Delete(ctx context.Context, id string) error {
	_, ok := m.collections[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(m.collections, id)
	return nil
}

func (m *CollectionRepository) FindByUser(ctx context.Context, userId string) ([]*domain.Collection, error) {
	result := make([]*domain.Collection, 0)
	for _, collection := range m.collections {
		if collection.OwnerId == userId {
			result = append(result, collection)
		}
	}
	return result, nil
}

func (m *CollectionRepository) AddParticipant(ctx context.Context, collectionId string, userId string) error {
	return nil
}

func (m *CollectionRepository) DeleteParticipant(ctx context.Context, collectionId string, userId string) error {
	return nil
}
