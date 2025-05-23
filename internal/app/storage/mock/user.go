package mock

import (
	"context"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
)

type UserRepository struct {
	users map[string]*domain.User
}

func NewUserRepository() domain.UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	result := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (m *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.users[id]
	return ok, nil
}

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.Id]; ok {
		return domain.ErrDuplicateEntity
	}
	m.users[user.Id] = user
	return nil
}

func (m *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, ok := m.users[user.Id]
	if !ok {
		return domain.ErrNotFound
	}
	m.users[user.Id] = user
	return nil
}

func (m *UserRepository) Delete(ctx context.Context, id string) error {
	_, ok := m.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *UserRepository) FindByMail(ctx context.Context, mail string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Mail == mail {
			return user, nil
		}
	}
	return nil, domain.ErrNotFound
}
