package database

import "github.com/waanvieira/api-users/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(user *entity.User) error
	FindByID(email string) (*entity.User, error)
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	Update(user *entity.User) error
	Delete(id string) error
}
