package service

import (
	"context"
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository domain.UserRepository
}

func (s *UserService) FindAll(ctx context.Context) ([]*domain.User, error) {
	return s.UserRepository.FindAll(ctx)
}

func (s *UserService) FindById(ctx context.Context, id string) (*domain.User, error) {
	return s.UserRepository.FindById(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *domain.User) (string, error) {
	user.Id = uuid.New().String()
	user.Password, _ = s.hashPassword(user.Password)
	return user.Id, s.UserRepository.Create(ctx, user)
}

func (s *UserService) FindOrCreate(ctx context.Context, user *domain.User) (*domain.User, error) {

	var err error
	var newId string

	if len(user.Id) == 0 {
		newId, err = s.Create(ctx, user)
	}
	if exists, err := s.UserRepository.Exists(ctx, user.Id); err == nil && !exists {
		newId, err = s.Create(ctx, user)
	}

	if err != nil {
		return nil, err
	}
	return s.FindById(ctx, newId)
}

func (s *UserService) UpdateUserDetails(ctx context.Context, user *domain.User) error {

	currentUser, err := s.FindById(ctx, user.Id)
	if err != nil {
		return err
	}

	var newUser *domain.User = &domain.User{
		Id:       user.Id,
		Password: currentUser.Password,
		Mail:     user.Mail,
	}

	return s.UserRepository.Update(ctx, newUser)
}

func (s *UserService) UpdatePassword(ctx context.Context, userId string, oldPassword string, newPassword string) error {

	currentUser, err := s.FindById(ctx, userId)
	if err != nil {
		return err
	}

	if s.verifyPassword(oldPassword, currentUser.Password) {

		currentUser.Password, _ = s.hashPassword(newPassword)
		return s.UserRepository.Update(ctx, currentUser)
	}

	return errors.New("invalid password")
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.UserRepository.Delete(ctx, id)
}

func (s *UserService) FindByMail(ctx context.Context, mail string) (*domain.User, error) {
	return s.UserRepository.FindByMail(ctx, mail)
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *UserService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
