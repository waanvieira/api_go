package database

import "github.com/waanvieira/api-users/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindByID(email string) (*entity.Product, error)
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
