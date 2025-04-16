package domain

import (
	"context"
)

type User struct {
	Id       string `gorm:"type:uuid;primary_key"`
	Mail     string `gorm:"unique"`
	Password string `json:"-"`
}

func NewUser(mail, password string) *User {
	return &User{
		Id:       "",
		Mail:     mail,
		Password: password,
	}
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	FindByMail(ctx context.Context, mail string) (*User, error)
}
